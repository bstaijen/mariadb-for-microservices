app.controller('ProfileController', function ($scope, $window, LocalStorage, ApiService) {

    if (!LocalStorage.hasToken() && !LocalStorage.hasUser()) {
        console.info('Token or User is missing from LocalStorage.');
        $window.location.href = '#/login';
        return;
    }


    $scope.errorMessage = "";
    $scope.successMessage = "";

    $scope.token = LocalStorage.getToken();
    $scope.user = LocalStorage.getUser();

    $scope.createdAt = displayMoment($scope.user.createdAt);

    function displayMoment(createdAt) {
        return moment(createdAt, "YYYY-MM-DD HH:mm:ss").fromNow();
    }

    $scope.rdonly = true;

    $scope.toggleEdit = function () {
        $scope.rdonly = !$scope.rdonly;
    };

    $scope.saveChanges = function () {
        $scope.errorMessage = "";
        $scope.successMessage = "";

        $scope.rdonly = !$scope.rdonly;

        ApiService.updateUser($scope.user).then(
            function (data) {
                LocalStorage.setUser(data);
                $scope.successMessage = "Profile succesfully updated."
            }, function (error) {
                console.error(error);
                $scope.errorMessage = "Error while updating profile."
            }
        );
    };

});