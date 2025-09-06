package services

import (
	"company-ai-training/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(email, name, department, position, employeeID string, startDate time.Time) (*models.User, error) {
	user := &models.User{
		ID:         uuid.New(),
		Email:      email,
		Name:       name,
		Department: department,
		Position:   position,
		EmployeeID: employeeID,
		StartDate:  startDate,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	updates["updated_at"] = time.Now()
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(userID uuid.UUID) error {
	return s.db.Delete(&models.User{}, "id = ?", userID).Error
}

// GetUserContext builds context string for AI about user
func (s *UserService) GetUserContext(user *models.User) (string, error) {
	now := time.Now()
	tenure := now.Sub(user.StartDate)
	yearsWorked := int(tenure.Hours() / (24 * 365))
	monthsWorked := int(tenure.Hours() / (24 * 30))

	context := fmt.Sprintf(`THÔNG TIN NHÂN VIÊN:
- Tên: %s
- Email: %s
- Mã nhân viên: %s
- Phòng ban: %s
- Chức vụ: Director
- Ngày bắt đầu làm việc: %s
- Thâm niên: %d năm %d tháng`,
		user.Name,
		user.Email,
		user.EmployeeID,
		user.Department,
		user.StartDate.Format("02/01/2006"),
		yearsWorked,
		monthsWorked%12,
	)

	return context, nil
}
