app.controller('DashboardController', function ($scope, $window, LocalStorage, ApiService) {

    if (!LocalStorage.hasToken() && !LocalStorage.hasUser()) {
        console.info('Token or User is missing from LocalStorage.');
        $window.location.href = '#/login';
        return;
    }

    $scope.token = LocalStorage.getToken();
    $scope.user = LocalStorage.getUser();
    $scope.error = "";

    console.log($scope.user);

    $scope.submit = function () {
        if (!$scope.title) {
            $scope.error = "title is mandatory";
            return;
        }

        ApiService.upload($scope.user.id, $scope.file, $scope.title).then(
            function (resp) {
                console.log('Success ' + resp.config.data.file.name + 'uploaded. Response: ' + resp.data);
                $scope.file = null;
            }, function (resp) {
                console.log('Error status: ' + resp.status);
            }, function (evt) {
                var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
                console.log('progress: ' + progressPercentage + '% ' + evt.config.data.file.name);
            }
        );
    }

    $scope.createdAt = displayMoment($scope.user.createdAt);

    function displayMoment(createdAt) {
        return moment(createdAt, "YYYY-MM-DD HH:mm:ss").fromNow();
    }

    $scope.rdonly = true;

    $scope.toggleEdit = function () {
        $scope.rdonly = !$scope.rdonly;
    };

    $scope.saveChanges = function () {
        $scope.rdonly = !$scope.rdonly;
        // TODO : save the changes in the backend
        ApiService.updateUser($scope.user).then(
            function (data) {
                LocalStorage.setUser(data);
            }, function (error) {
                console.error(error);
            }
        );
    };

});