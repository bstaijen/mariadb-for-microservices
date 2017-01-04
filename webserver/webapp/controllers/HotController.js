app.controller('HotController', function ($scope, ApiService) {
    $scope.photos = [];

    ApiService.hot().then(
        function (result) {
            console.log(result);

            $scope.photos = result;

        }, function (error) {
            console.error(error);
        }
    );


});