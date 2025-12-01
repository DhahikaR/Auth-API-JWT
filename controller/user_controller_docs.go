package controller

// CreateUser godoc
// @Summary Create new user (Admin only)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body web.UserCreateRequest true "Create user"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Router /users [post]
func (UserControllerImpl) CreateUserDocs() {}

// FindAllUsers godoc
// @Summary Get all users (Admin only)
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} web.WebResponse
// @Router /users [get]
func (UserControllerImpl) FindAllDocs() {}

// FindUserById godoc
// @Summary Get user by ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} web.WebResponse
// @Router /users/{userId} [get]
func (UserControllerImpl) FindByIdDocs() {}

// UpdateUser godoc
// @Summary Update user (Admin only)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param request body web.UserUpdateRequest true "Update user"
// @Success 200 {object} web.WebResponse
// @Router /users/{userId} [put]
func (UserControllerImpl) UpdateUserDocs() {}

// DeleteUser godoc
// @Summary Delete user (Admin only)
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} web.WebResponse
// @Router /users/{userId} [delete]
func (UserControllerImpl) DeleteUserDocs() {}

// Me godoc
// @Summary Get user profile (Me)
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} web.WebResponse
// @Router /users/me [get]
func (UserControllerImpl) MeDocs() {}

// UpdateMe godoc
// @Summary Update user profile (Me)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body web.UserUpdateRequest true "Update user"
// @Success 200 {object} web.WebResponse
// @Router /users/me [put]
func (UserControllerImpl) UpdateMeDocs() {}
