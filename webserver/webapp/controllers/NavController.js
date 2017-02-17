app.controller('NavController', function ($scope, $location, $route, LocalStorage) {
    $scope.$route = $route;


    $scope.isLoggedIn = function() {
        return LocalStorage.hasToken();
    };

    $scope.openUploadWindow = function() {
        $location.path('/upload');
    }

});