app.directive('post', function (LocalStorage, ApiService, $uibModal) {
    return {
        restrict: 'E',
        scope: {
            photo: '='
        },
        templateUrl: 'directives/postView.html',
        controller: function ($scope) {


            $scope.lastCommentsLoaded = true;
            $scope.showMessages = false;
            $scope.page = 0;

            if ($scope.photo.comments.length == 10) {
                $scope.lastCommentsLoaded = false;
                // set page to 1
                $scope.page = 1;
            }

            $scope.loadComments = function () {

                // make click listener.
                // increase page ++1
                $scope.page = $scope.page + 1;

                var page = $scope.page;
                var itemsPerPage = 10;
                var offset = (page - 1) * itemsPerPage + 1;

                // get comments request.
                ApiService.getComments($scope.photo.id, offset, 10).then(
                    function (response) {
                        console.info(response);

                        //  : if loaded comments are size of 10 then lastCommentsLoaded stays false, else true
                        if (response.length == 10) {
                            $scope.lastCommentsLoaded = false;
                        } else {
                            $scope.lastCommentsLoaded = true;
                        }

                        $scope.photo.comments = $.merge($scope.photo.comments, response);
                    }, function (error) {
                        console.error(error);
                    }
                );


            };

            $scope.toggleMessages = function () {
                var usr = LocalStorage.getUser();
                if (!usr) {
                    showSignInRequiredModal();
                    return;
                }

                $scope.showMessages = !$scope.showMessages;
            };

            $scope.buildImageUrl = function (filename) {
                return ApiService.urlbuilder.photo('/images/' + filename);
            };

            $scope.displayMoment = function (createdAt) {
                return moment(createdAt, "YYYY-MM-DD HH:mm:ss").fromNow();
            };

            $scope.upvote = function (photo) {
                console.log('upvote button clicked');
                if (photo.upvote) {
                    console.log('You already upvoted');
                    return;
                }
                var usr = LocalStorage.getUser();
                if (!usr) {
                    showSignInRequiredModal();
                    return;
                }

                ApiService.upvote(usr.id, photo.id).then(function (resp) {
                    // If this is the first vote from this user then increase the number of total votes.
                    if (!photo.downvote && !photo.upvote) {
                        photo.totalVotes++;
                        photo.upvote_count++;
                    }

                    // If the previous vote was a upvote, then recalculate the counts
                    if (photo.downvote) {
                        photo.downvote_count--;
                        photo.upvote_count++;
                    }

                    // Let angular know which button to show on page
                    photo.downvote = false;
                    photo.upvote = true;

                }, function (error) {
                    console.error(error);
                });

            };

            $scope.downvote = function (photo) {
                console.log('downvote button clicked');
                if (photo.downvote) {
                    console.log('You already downvoted');
                    return;
                }
                var usr = LocalStorage.getUser();
                if (!usr) {
                    showSignInRequiredModal();
                    return;
                }

                ApiService.downvote(usr.id, photo.id).then(function (resp) {
                    // If this is the first vote from this user then increase the number of total votes.
                    if (!photo.downvote && !photo.upvote) {
                        photo.totalVotes++;
                        photo.downvote_count++;
                    }

                    // If the previous vote was a upvote, then recalculate the counts
                    if (photo.upvote) {
                        photo.upvote_count--;
                        photo.downvote_count++;
                    }

                    // Let angular know which button to show on page
                    photo.upvote = false;
                    photo.downvote = true;

                }, function (error) {
                    console.error(error);
                });
            };

            $scope.comment = function (photo) {
                var usr = LocalStorage.getUser();
                if (!usr) {
                    showSignInRequiredModal();
                    return;
                }

                ApiService.comment(usr.id, photo.id, $scope.comment_text).then(
                    function (response) {
                        console.info(response);
                        $scope.comment_text = '';
                        $scope.photo.comment_count++;
                        if ($scope.photo && $scope.photo.comments) {
                            $scope.photo.comments.push(response);
                        }
                    },
                    function (error) {
                        console.error(error);
                    }
                );
            }
            $scope.calcPerc = function (photo) {

                var down = photo.downvote_count < 1 ? 0 : photo.downvote_count * 100;
                var up = photo.upvote_count < 1 ? 0 : photo.upvote_count * 100;

                if (down < 1 && up < 1) {
                    return '-';
                } else if (down < 1) {
                    return '100 %'
                } else if (up < 1) {
                    return '0 %'
                } else {
                    return Math.floor((up / (up + down) * 100)) + " %"
                }

            };

            function showSignInRequiredModal() {
                var modalInstance = $uibModal.open({
                    animation: true,
                    ariaLabelledBy: 'modal-title',
                    ariaDescribedBy: 'modal-body',
                    templateUrl: 'view/modal.html',
                    controller: function ($uibModalInstance) {
                        var $ctrl = this;
                        $ctrl.cancel = function () {
                            $uibModalInstance.dismiss('cancel');
                        };
                    },
                    controllerAs: '$ctrl',
                    size: 'sm'
                });

                // To disable error message in console:
                modalInstance.closed.then(function () {
                }, function () {
                });
                // To disable error message in console:
                modalInstance.result.then(function () {
                }, function () {
                });
            }
        }
    }
});