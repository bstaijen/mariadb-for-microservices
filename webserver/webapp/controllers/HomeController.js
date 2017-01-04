app.controller('HomeController', function ($scope, ApiService) {
    $scope.photos = [];

    ApiService.incoming().then(
        function (result) {
            console.log(result);

            $scope.photos = result;

        }, function (error) {
            console.error(error);
        }
    );
});