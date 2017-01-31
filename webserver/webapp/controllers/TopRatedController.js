app.controller('TopRatedController', function ($scope, ApiService) {
    // collection of photos.
    $scope.photos = [];

    // keeps track of how many pages of photos has been loaded.
    $scope.page = 1;

    // boolean for keeping track if new photos can be loaded
    $scope.hasMorePictures = true;

    // contains error messages
    $scope.errors = [];

    // load photos
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
                console.error(error);

                if (error.data == null && error.status == -1) {
                    $scope.errors.push("Error connecting to API. Maybe resource is offline?");
                }
            }
        );
    }

    // function for loading the next 10 photos.
    $scope.loadNext = function () {
        $scope.page = $scope.page + 1;
        getPhotos();
    }
});