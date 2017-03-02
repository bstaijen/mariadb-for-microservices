app.controller('UserController', function ($scope, $window, LocalStorage, ApiService, $location) {

    if (!LocalStorage.hasToken() && !LocalStorage.hasUser()) {
        console.info('Token and/or User is missing from LocalStorage.');
        $window.location.href = '#/login';
        return;
    }

    $scope.username = LocalStorage.getUser().username;

    // keeps track of how many pages of photos has been loaded.
    $scope.photo_pages = 1;
    $scope.photos = [];

    $scope.voted_photos = [];

    $scope.comment_tupel = [];

    $scope.commentsTabVisible = false;
    $scope.votesTabVisible = false;
    $scope.photosTabVisible = true;

    $scope.showLoading = false;

    // To load the picture on startup
    getPhotos();

    $scope.showCommentsTab = function () {
        $scope.commentsTabVisible = true;
        $scope.votesTabVisible = false;
        $scope.photosTabVisible = false;
        getComments();
    };

    $scope.showVotesTab = function () {
        $scope.commentsTabVisible = false;
        $scope.votesTabVisible = true;
        $scope.photosTabVisible = false;
        getVotes();
    };

    $scope.showPhotosTab = function () {
        $scope.commentsTabVisible = false;
        $scope.votesTabVisible = false;
        $scope.photosTabVisible = true;
        getPhotos();
    };

    $scope.buildImageUrl = function (filename) {
        return ApiService.urlbuilder.photo('/images/' + filename);
    };

    $scope.displayMoment = function (createdAt) {
        return moment(createdAt, "YYYY-MM-DD HH:mm:ss").fromNow();
    };
    $scope.calcPerc = function (photo) {

        var down = photo.downvote_count < 1 ? 0 : photo.downvote_count * 100;
        var up = photo.upvote_count < 1 ? 0 : photo.upvote_count * 100;

        if (down < 1 && up < 1) {
            return '0%';
        } else if (down < 1) {
            return '100 %'
        } else if (up < 1) {
            return '0 %'
        } else {
            return Math.floor((up / (up + down) * 100)) + " %"
        }

    };

    function getComments() {
        $scope.showLoading = true;
        ApiService.getCommentsFromUser().then(
            function (response) {
                console.info(response);
                if (response.result) {
                    $scope.comment_tupel = response.result;
                }
                $scope.showLoading = false;
            }, function (error) {
                console.error(error);
                $scope.showLoading = false;
            }
        )
    }


    function getPhotos() {
        $scope.showLoading = true;

        var page = $scope.photo_pages;
        var itemsPerPage = 10;
        var offset = (page - 1) * itemsPerPage;

        var user_id = LocalStorage.getUser().id;

        ApiService.getPhotosForUser(user_id, offset, itemsPerPage).then(
            function (response) {
                console.info(response);
                if (response) {
                    $scope.photos = response;
                }
                $scope.showLoading = false;
            }, function (error) {
                console.error(error);
                $scope.showLoading = false;
            }
        );
    }

    function getVotes() {
        $scope.showLoading = true;
        ApiService.getTheVotesFromUser().then(
            function (response) {
                console.info(response);
                if (response.result) {
                    $scope.voted_photos = response.result;
                }
                $scope.showLoading = false;
            }, function (error) {
                console.error(error);
                $scope.showLoading = false;
            }
        )
    }

    $scope.deleteComment = function(commentID) {
        $scope.showLoading = true;
        ApiService.deleteComment(commentID).then(
            function(response){
                console.info(response);

                angular.forEach($scope.comment_tupel, function (tupel, index) {
                    if (tupel.comment.id === commentID) {
                        $scope.comment_tupel.splice(index, 1);
                    }
                });
                $scope.showLoading = false;
            },
            function(error){
                console.error(error);
                $scope.showLoading = false;
            }
        );
    };

    $scope.deletePhoto = function(photoID) {
        $scope.showLoading = true;
        ApiService.deletePhoto(photoID).then(
            function(response){
                console.info(response);

                angular.forEach($scope.photos, function (photo, index) {
                    if (photo.id === photoID) {
                        $scope.photos.splice(index, 1);
                    }
                })
                $scope.showLoading = false;
            },
            function(error){
                console.error(error);
                $scope.showLoading = false;
            }
        );
    };

    $scope.openImage = function(id) {
        if (!id) {
            console.warn('openImage: id is empty');
            return;
        }
        $location.path('/photo/' + id);
    }

});