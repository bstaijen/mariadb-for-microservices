app.controller('TopRatedController', function ($scope, ApiService) {
    $scope.photos = [];

    ApiService.toprated().then(
        function (result) {
            console.log(result);

            $scope.photos = result;

        }, function (error) {
            console.error(error);
        }
    );


});