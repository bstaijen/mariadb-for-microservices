app.controller('TopRatedController', function ($scope, ApiService) {
    $scope.photos = [];
    $scope.page = 1;
    $scope.hasMorePictures = true;

    getPhotos();
    function getPhotos() {
        var page = $scope.page;
        var itemsPerPage = 10;
        var offset = (page - 1) * itemsPerPage;

        ApiService.toprated(offset, itemsPerPage).then(
            function (result) {
                console.log(result);

                $scope.hasMorePictures = (result.length == itemsPerPage);

                $scope.photos = $.merge($scope.photos, result);

            }, function (error) {
                // TODO
                console.error(error);
            }
        );
    }

    $scope.loadNext = function () {
        $scope.page = $scope.page + 1;
        getPhotos();
    }
});