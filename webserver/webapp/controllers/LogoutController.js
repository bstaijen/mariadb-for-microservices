app.controller("LogoutController", function ($window, LocalStorage) {
    LocalStorage.removeToken();
    LocalStorage.removeUser();

    $window.location.href = '#/';

});