package service

import (
	// "strings" // No longer used directly
	"testing"
	"time"
    "errors"
    // "fmt" // Not used directly, log or t.Logf is used

	"github.com/russian-steam/auth-service/internal/domain"
	"github.com/google/uuid"
)

const testVerificationCodeExpiry = 15 * time.Minute // Align with service const if possible

func TestVerificationService_GenerateEmailVerificationCode(t *testing.T) {
	mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
	vs := NewVerificationService(mockVCRepo, mockUserRepo)

	userID := "user_verify_gen"
	email := "verify_gen@example.com"

	vc, rawCode, err := vs.GenerateEmailVerificationCode(userID, email)
	if err != nil {
		t.Fatalf("GenerateEmailVerificationCode() error = %v", err)
	}
	if vc == nil || rawCode == "" {
		t.Fatalf("GenerateEmailVerificationCode() returned nil vc or empty rawCode")
	}
    if len(rawCode) != verificationCodeLength {
         t.Errorf("GenerateEmailVerificationCode() raw code length got = %d, want %d", len(rawCode), verificationCodeLength)
    }
	if vc.UserID != userID || vc.Target != email || vc.Type != domain.VerificationTypeEmail {
		t.Errorf("GenerateEmailVerificationCode() vc data mismatch")
	}
    if time.Until(vc.ExpiresAt).Minutes() < (testVerificationCodeExpiry.Minutes() -1) || time.Until(vc.ExpiresAt) > testVerificationCodeExpiry {
        t.Errorf("GenerateEmailVerificationCode() expiry time not set correctly, remaining: %v, expected around: %v", time.Until(vc.ExpiresAt), testVerificationCodeExpiry)
    }

    // Check if stored in mock (by hash)
    hashedCode := vs.hashVerificationCode(rawCode)
    storedVC, ok := mockVCRepo.codes[hashedCode]
    if !ok {
        t.Errorf("GenerateEmailVerificationCode() code not stored in mock repo by hash %s", hashedCode)
    } else if storedVC.ID != vc.ID {
         t.Errorf("GenerateEmailVerificationCode() stored VC ID mismatch")
    }
}

func TestVerificationService_VerifyEmailCode_Success(t *testing.T) {
    mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
    vs := NewVerificationService(mockVCRepo, mockUserRepo)

    userID := uuid.NewString()
    emailToVerify := "verify_success@example.com"

    // Setup user that needs verification
    userToVerify := &domain.User{ID: userID, Username: "verifyMe", Email: emailToVerify, Status: domain.StatusPendingVerification, CreatedAt: time.Now(), UpdatedAt: time.Now()}
    mockUserRepo.Create(userToVerify)

    // Generate a code first using the service itself to ensure consistency
    generatedVC, rawCodeToTest, errGen := vs.GenerateEmailVerificationCode(userID, emailToVerify)
    if errGen != nil {
        t.Fatalf("Setup: GenerateEmailVerificationCode() failed: %v", errGen)
    }

    verifiedUser, err := vs.VerifyEmailCode(rawCodeToTest, emailToVerify)
    if err != nil {
        t.Fatalf("VerifyEmailCode() error = %v", err)
    }
    if verifiedUser == nil {
        t.Fatalf("VerifyEmailCode() returned nil user")
    }
    if verifiedUser.ID != userID {
        t.Errorf("VerifyEmailCode() returned user ID %s, want %s", verifiedUser.ID, userID)
    }
    if verifiedUser.Status != domain.StatusActive {
        t.Errorf("VerifyEmailCode() user status got = %s, want %s", verifiedUser.Status, domain.StatusActive)
    }
    if verifiedUser.EmailVerifiedAt == nil || verifiedUser.EmailVerifiedAt.IsZero() {
         t.Errorf("VerifyEmailCode() user EmailVerifiedAt not set")
    }

    // Check if code was deleted from mock
    hashedCode := vs.hashVerificationCode(rawCodeToTest)
    if _, ok := mockVCRepo.codes[hashedCode]; ok {
        t.Errorf("VerifyEmailCode() code with hash %s not deleted from mock repo after use", hashedCode)
    }
    // Also check by original VC ID if mock supports that (it should be deleted by ID internally)
    var foundByID bool
    for _, code := range mockVCRepo.codes { // Iterate as Delete might be by ID but storage by hash
        if code.ID == generatedVC.ID {
            foundByID = true
            break
        }
    }
    if foundByID {
         t.Errorf("VerifyEmailCode() code with ID %s not deleted from mock repo after use", generatedVC.ID)
    }
}

func TestVerificationService_VerifyEmailCode_UserAlreadyActive(t *testing.T) {
    mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
    vs := NewVerificationService(mockVCRepo, mockUserRepo)

    userID := uuid.NewString()
    email := "already_active@example.com"
    now := time.Now().UTC()

    userAlreadyActive := &domain.User{ID: userID, Email: email, Status: domain.StatusActive, EmailVerifiedAt: &now}
    mockUserRepo.Create(userAlreadyActive)

    _, rawCode, _ := vs.GenerateEmailVerificationCode(userID, email) // Generate a code

    verifiedUser, err := vs.VerifyEmailCode(rawCode, email)
    if err != nil {
        t.Fatalf("VerifyEmailCode() for already active user error = %v", err)
    }
    if verifiedUser.Status != domain.StatusActive { // Should remain active
        t.Errorf("VerifyEmailCode() user status changed from active, got %s", verifiedUser.Status)
    }
     if verifiedUser.EmailVerifiedAt == nil || !verifiedUser.EmailVerifiedAt.Equal(now) { // Should remain as original verification time
        t.Errorf("VerifyEmailCode() user EmailVerifiedAt changed or unset, got %v, want %v", verifiedUser.EmailVerifiedAt, now)
    }
}


func TestVerificationService_VerifyEmailCode_NotFound(t *testing.T) {
    mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
    vs := NewVerificationService(mockVCRepo, mockUserRepo)

    _, err := vs.VerifyEmailCode("nonexistentcode", "test@example.com")
    if !errors.Is(err, ErrVerificationCodeNotFound) {
        t.Errorf("VerifyEmailCode() with non-existent code, error = %v, want %v", err, ErrVerificationCodeNotFound)
    }
}

func TestVerificationService_VerifyEmailCode_Expired(t *testing.T) {
    mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
    vs := NewVerificationService(mockVCRepo, mockUserRepo)

    userID := "user_code_expired"
    email := "expired@example.com"
    rawCode := "123456"
    hashedCode := vs.hashVerificationCode(rawCode)
    expiredVCID := uuid.NewString()

    mockVCRepo.Create(&domain.VerificationCode{
        ID: expiredVCID, UserID: userID, Type: domain.VerificationTypeEmail, CodeHash: hashedCode, Target: email,
        ExpiresAt: time.Now().UTC().Add(-1 * time.Minute), // Expired
        CreatedAt: time.Now().UTC().Add(-16 * time.Minute),
    })

    // User to be updated (or attempted)
    mockUserRepo.Create(&domain.User{ID: userID, Email: email, Status: domain.StatusPendingVerification})

    _, err := vs.VerifyEmailCode(rawCode, email)
    if !errors.Is(err, ErrVerificationCodeExpired) {
        t.Errorf("VerifyEmailCode() with expired code, error = %v, want %v", err, ErrVerificationCodeExpired)
    }
    // Check if expired code was deleted from mock
    var foundByID bool
    for _, code := range mockVCRepo.codes {
        if code.ID == expiredVCID {
            foundByID = true
            break
        }
    }
    if foundByID {
         t.Errorf("VerifyEmailCode() expired code with ID %s not deleted from mock repo", expiredVCID)
    }
}

func TestVerificationService_VerifyEmailCode_TargetMismatch(t *testing.T) {
    mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
    vs := NewVerificationService(mockVCRepo, mockUserRepo)

    userID := "user_target_mismatch"
    correctEmail := "correct@example.com"
    wrongEmailAttempt := "wrong@example.com"
    rawCode := "abcdef"
    hashedCode := vs.hashVerificationCode(rawCode)

    mockVCRepo.Create(&domain.VerificationCode{
        ID: uuid.NewString(), UserID: userID, Type: domain.VerificationTypeEmail, CodeHash: hashedCode, Target: correctEmail,
        ExpiresAt: time.Now().UTC().Add(testVerificationCodeExpiry), CreatedAt: time.Now().UTC(),
    })

    // User exists for verification service to potentially interact with if target matched
    mockUserRepo.Create(&domain.User{ID: userID, Email: correctEmail, Status: domain.StatusPendingVerification})

    _, err := vs.VerifyEmailCode(rawCode, wrongEmailAttempt) // Attempt to verify with wrong email
    if !errors.Is(err, ErrVerificationCodeInvalid) {
        t.Errorf("VerifyEmailCode() with wrong target email, error = '%v', want an error wrapping '%v'", err, ErrVerificationCodeInvalid)
    }
    t.Logf("VerifyEmailCode() with wrong target email, got error: %v (expected ErrVerificationCodeInvalid or similar)", err)
}

func TestVerificationService_GenerateEmailVerificationCode_DeletesOldCodes(t *testing.T) {
    mockVCRepo := NewMockVerificationCodeRepository()
    mockUserRepo := NewMockUserRepository()
    vs := NewVerificationService(mockVCRepo, mockUserRepo)

    userID := "user_gen_delete_old"
    email := "delete_old@example.com"

    // Create an old code
    oldRawCode := "old123"
    oldHashedCode := vs.hashVerificationCode(oldRawCode)
    mockVCRepo.Create(&domain.VerificationCode{
        ID: uuid.NewString(), UserID: userID, Type: domain.VerificationTypeEmail, CodeHash: oldHashedCode, Target: email,
        ExpiresAt: time.Now().UTC().Add(testVerificationCodeExpiry), CreatedAt: time.Now().UTC().Add(-time.Hour),
    })

    // Generate new code
    _, newRawCode, err := vs.GenerateEmailVerificationCode(userID, email)
    if err != nil {
        t.Fatalf("GenerateEmailVerificationCode() error = %v", err)
    }

    // Check old code is deleted
    if _, ok := mockVCRepo.codes[oldHashedCode]; ok {
        t.Errorf("GenerateEmailVerificationCode() did not delete old code with hash %s", oldHashedCode)
    }
    // Check new code is present
    newHashedCode := vs.hashVerificationCode(newRawCode)
    if _, ok := mockVCRepo.codes[newHashedCode]; !ok {
        t.Errorf("GenerateEmailVerificationCode() new code with hash %s not found", newHashedCode)
    }
}
