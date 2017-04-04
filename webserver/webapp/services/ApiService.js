app.factory('ApiService', function ($http, $q, LocalStorage, Upload) {

    var photo_url = '.';
    var authentication_url = '.';
    var profile_url = '.';
    var vote_url = '.';
    var comment_url = '.';
    var settings = {
        debug: true
    };

    function composeAuthenticationUrl(url) {
        return authentication_url + url;
    }

    function composeProfileUrl(url) {
        return profile_url + url;
    }

    function composePhotoUrl(url) {
        return photo_url + url;
    }

    function composeVoteUrl(url) {
        return vote_url + url;
    }

    function composeCommentUrl(url) {
        return comment_url + url;
    }

    function call() {
        // TODO refactor that all calls go through this function
    }

    function get(url) {
        if (settings.debug) {
            console.log('API: GET ' + url);
        }

        var promise = $q.defer();
        $http.get(url, {
            headers: {
                // TODO : Authorization: getAccessToken()
            }
        }).then(function success(response) {
            promise.resolve(response.data);
        }, function error(response) { // Error callback

            if (settings.consoleDebug) {
                console.log(response);
            }

            // Handle error
            // TODO : errorHandler(response, promise);
            promise.reject(response);

        });
        return promise.promise;
    }

    function post(url, options) {

        if (settings.debug) {
            console.log('API: POST ' + url);
        }

        var promise = $q.defer();

        $http.post(url, options, {
            headers: {
                'Content-Type': 'application/json; charset=UTF-8'
            }
        }).then(function success(response) {
            promise.resolve(response.data);
        }, function error(response) { // Error callback

            if (settings.debug) {
                console.log(response);
            }
            // Handle error
            // TODO : errorHandler(response, promise);
            promise.reject(response);

        });

        return promise.promise;
    }

    function put(url, options) {

        console.info(options);

        if (settings.debug) {
            console.log('API: PUT ' + url);
        }

        var promise = $q.defer();

        $http.put(url, options, {
            headers: {
                'Content-Type': 'application/json; charset=UTF-8'
            }
        }).then(function success(response) {
            promise.resolve(response.data);
        }, function error(response) { // Error callback

            if (settings.debug) {
                console.log(response);
            }
            promise.reject(response);

        });

        return promise.promise;
    }

    return {
        login: function (username, password) {
            var part = '/token-auth';
            var url = composeAuthenticationUrl(part);
            return post(url, {
                username: username,
                password: password
            });
        },
        register: function (username, email, password) {
            var url = composeProfileUrl('/users');
            return post(url, {username: username, password: password, email: email});
        },
        upload: function (user_id, file, title) {
            var url = composePhotoUrl('/image/' + user_id + "?title=" + title + '&token=' + LocalStorage.getToken());
            return Upload.upload({
                url: url,
                data: {file: file}
            });
        },
        incoming: function (offset, rows) {
            if (LocalStorage.hasToken()) {
                return get(composePhotoUrl('/image/list?offset=' + offset + '&rows=' + rows + '&token=' + LocalStorage.getToken()));
            } else {
                return get(composePhotoUrl('/image/list?offset=' + offset + '&rows=' + rows));
            }
        },
        toprated: function (offset, rows) {
            if (LocalStorage.hasToken()) {
                return get(composePhotoUrl('/image/toprated?offset=' + offset + '&rows=' + rows + '&token=' + LocalStorage.getToken()));
            } else {
                return get(composePhotoUrl('/image/toprated?offset=' + offset + '&rows=' + rows));
            }
        },
        hot: function (offset, rows) {
            if (LocalStorage.hasToken()) {
                return get(composePhotoUrl('/image/hot?offset=' + offset + '&rows=' + rows + '&token=' + LocalStorage.getToken()));
            } else {
                return get(composePhotoUrl('/image/hot?offset=' + offset + '&rows=' + rows));
            }
        },
        updateUser: function (user) {
            var url = composeProfileUrl('/users?token=' + LocalStorage.getToken());
            return put(url, user);
        },
        upvote: function (user_id, photo_id) {
            var url = composeVoteUrl('/votes?token=' + LocalStorage.getToken());
            var options = {
                user_id: user_id,
                photo_id: photo_id,
                upvote: true
            };
            console.info(options);
            return post(url, options);
        },
        downvote: function (user_id, photo_id) {
            var url = composeVoteUrl('/votes?token=' + LocalStorage.getToken());
            var options = {
                user_id: user_id,
                photo_id: photo_id,
                downvote: true
            };
            console.info(options);
            return post(url, options);
        },
        comment: function (user_id, photo_id, comment) {
            var url = composeCommentUrl('/comments');
            var options = {
                user_id: user_id,
                photo_id: photo_id,
                comment: comment
            }
            console.info(options);
            return post(url, options);
        },
        getComments: function (photo_id, offset, nr_of_rows) {
            var url = composeCommentUrl('/comments?photoID=' + photo_id + "&offset=" + offset + "&rows=" + nr_of_rows);
            return get(url);
        },
        urlbuilder: {
            authenication: composeAuthenticationUrl,
            vote: composeVoteUrl,
            comment: composeCommentUrl,
            profile: composeProfileUrl,
            photo: composePhotoUrl
        },
        getPhotoByID: function(photoID) {
            if (LocalStorage.hasToken()) {
                var url = composePhotoUrl('/image/' + photoID + '?token=' + LocalStorage.getToken());
                return get(url);
            } else {
                var url = composePhotoUrl('/image/' + photoID);
                return get(url);
            }
        },
        getPhotosForUser: function(user_id, offset, nr_of_rows) {
            var url = composePhotoUrl('/image/' + user_id + "/list?offset=" + offset + "&rows=" + nr_of_rows);
            return get(url);
        },
        getTheVotesFromUser: function() {
            var url = composePhotoUrl('/votes?token=' + LocalStorage.getToken());
            return get(url);
        },
        getCommentsFromUser: function() {
            var url = composePhotoUrl('/comments/fromuser?token=' + LocalStorage.getToken());
            return get(url);
        },
        deleteComment: function(commentID) {
            var url = composeCommentUrl('/comments/' + commentID + '/delete?token=' + LocalStorage.getToken());
            return post(url, {})
        },
        deletePhoto: function(photoID) {
            var url = composePhotoUrl('/image/' + photoID + '/delete?token=' + LocalStorage.getToken());
            return post(url, {})
        }
    }
});