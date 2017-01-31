app.controller('UploadController', function ($scope, $window, LocalStorage, ApiService) {

    if (!LocalStorage.hasToken() && !LocalStorage.hasUser()) {
        console.info('Token and/or User is missing from LocalStorage.');
        $window.location.href = '#/login';
        return;
    }

    $scope.error = "";
    $scope.successMessage = "";

    $scope.submit = function () {
        if (!$scope.title) {
            $scope.error = "title is mandatory";
            return;
        }

        var user = LocalStorage.getUser();

        ApiService.upload(user.id, $scope.file, $scope.title).then(
            function (resp) {
                $scope.title = "";
                $scope.file = null;
                $scope.successMessage = "Image successfully uploaded.";
            }, function (resp) {
                console.log('Error status: ' + resp.status);
                $scope.error = "Error uploading file.";
            }, function (evt) {
                var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
                console.log('progress: ' + progressPercentage + '% ' + evt.config.data.file.name);
            }
        );
    };

});