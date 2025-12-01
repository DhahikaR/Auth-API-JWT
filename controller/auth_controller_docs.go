package controller

// Register godoc
// @Summary Register new user
// @Description Mendaftarkan user baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body web.AuthRegisterRequest true "Register payload"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Router /auth/register [post]
func (AuthControllerImpl) RegisterDocs() {}

// Login godoc
// @Summary Login user
// @Description Mengembalikan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body web.AuthLoginRequest true "Login payload"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Router /auth/login [post]
func (AuthControllerImpl) LoginDocs() {}
