LoginView     = require './AppView'

module.exports = class LoginAppsController extends KDViewController

  KD.registerAppClass this, name  : 'Login'


  constructor:(options = {}, data)->

    KD.registerSingleton 'loginController', this

    options.view    = new LoginView
      testPath      : "landing-login"
    options.appInfo =
      name          : "Login"

    super options, data


  handleQuery: ({ query }) ->

    loginView = @getOption 'view'
    loginView.setCustomData query


  prepareFinishRegistrationForm: (token) ->

    { JPasswordRecovery } = KD.remote.api
    JPasswordRecovery.fetchRegistrationDetails token, (err, details) =>
      view = @getView()
      if err
        KD.showError err
        view.animateToForm 'login'
        return

      view.finishRegistrationForm.setRegistrationDetails details
      view.setCustomDataToForm 'finishRegistration', recoveryToken: token
      view.animateToForm 'finishRegistration'


  headBannerShowInvitation: (invite) ->

    view = @getView()
    view.headBannerShowInvitation invite


  setStorageData: (key, value) ->

    @appStorage = KD.getSingleton('appStorageController').storage 'Login', '1.0'
    @appStorage.fetchStorage (storage) =>
      @appStorage.setValue key, value, (err) ->
        warn "Failed to set #{key} information"  if err

