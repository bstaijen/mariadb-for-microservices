app.controller('NavController', function ($scope, $route, LocalStorage) {
    $scope.$route = $route;


    $scope.isLoggedIn = function() {
        return LocalStorage.hasToken();
    };

});