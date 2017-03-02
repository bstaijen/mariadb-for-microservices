app.controller('PhotoDetailsController', function ($routeParams, $scope, $window, LocalStorage, ApiService, $compile) {

    $scope.photoObject = {};

    var photoID = $routeParams.id;
    if (!photoID || photoID < 1) {
        console.warn("Photo ID is wrong.");
        console.warn("PhotoID: " + photoID);
        return
    }

    console.info(photoID);

    ApiService.getPhotoByID(photoID).then(
        function (response) {
            console.info(response);
            if (response) {
                $scope.photoObject = response;


                var $element = angular.element(document.querySelector('#photoDetails'));
                $element.children().remove();

                var photoElement = $compile('<post photo="photoObject"></post>')($scope);
                $element.append(photoElement);
            }
        }, function (error) {
            console.error(error);
        }
    );

});