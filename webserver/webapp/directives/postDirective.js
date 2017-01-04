app.directive('post', function(LocalStorage, ApiService){
    return {
        restrict : 'E',
        scope: {
            photo: '='
        },
        templateUrl : 'directives/postView.html',
        controller: function($scope) {

            $scope.showMessages = false;

            $scope.toggleMessages = function() {
                $scope.showMessages = !$scope.showMessages;
            };

            $scope.buildImageUrl = function (filename) {
                return ApiService.urlbuilder.photo('/images/' + filename);
            };

            $scope.displayMoment = function (createdAt) {
                return moment(createdAt, "YYYY-MM-DD HH:mm:ss").fromNow();
            };


            $scope.upvote = function (photo) {
                console.log('upvote');
                if (photo.upvote) {
                    console.log('photo.upvote==true');
                    return;
                }
                var usr = LocalStorage.getUser();
                if (!usr) {
                    // TODO : maybe redirect to login?
                    // Or show a message with : 'you need to signup buddy'
                    console.info('you need to singin buddy')
                    return;
                }

                ApiService.upvote(usr.id, photo.id).then(function (resp) {
                    if (!photo.downvote && !photo.upvote) {
                        photo.totalVotes++;
                    }
                    photo.downvote = false;
                    photo.upvote = true;
                }, function (error) {
                    console.error(error);
                });

            };
            $scope.downvote = function(photo) {
                console.log('downvote');
                if (photo.downvote) {
                    console.log('photo.downvote==true');
                    return;
                }
                var usr = LocalStorage.getUser();
                if (!usr) {
                    console.info('you need to singin buddy')
                    // TODO : maybe redirect to login?
                    // Or show a message with : 'you need to signup buddy'
                    return;
                }

                ApiService.downvote(usr.id, photo.id).then(function (resp) {
                    if (!photo.downvote && !photo.upvote) {
                        photo.totalVotes++;
                    }
                    photo.upvote = false;
                    photo.downvote = true;
                }, function (error) {
                    console.error(error);
                });
            };
            $scope.comment = function(photo) {
                var usr = LocalStorage.getUser();
                if (!usr) {
                    // TODO : maybe redirect to login?
                    // Or show a message with : 'you need to signup buddy'
                    return;
                }

                // TODO : get ng-model textarea

                ApiService.comment(usr.id, photo.id, $scope.comment_text).then(
                    function (response) {
                        console.info(response);
                        $scope.comment_text = '';
                        if ($scope.photo && $scope.photo.comments) {
                            $scope.photo.comments.push(response);
                        }
                    },
                    function (error) {
                        console.error(error);
                    }
                );
            }
        }
    }
});