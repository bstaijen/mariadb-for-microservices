app.controller('RegisterController', function ($scope, $window, ApiService, LocalStorage) {
    $scope.errorMessages = [];
    $scope.successMessages = [];

    $scope.register = function () {
        emptyMessages();

        ApiService.register($scope.username, $scope.email, $scope.password).then(
            function (data) {
                console.log(data);
                $scope.successMessages.push("Registration successful");
                emptyFields();

                LocalStorage.setToken(data.token);

                if (data.user) {
                    //console.info(data.user);
                    LocalStorage.setUser(data.user);
                }


                $window.location.href = '#/';

            }, function (response) {
                console.log(response);
                if (response && response.data) {
                    var data = response.data;

                    if (data && data.message) {
                        $scope.errorMessages.push(data.message);
                        return;
                    }
                }
                $scope.errorMessages.push("An error occured. Please try again.");
            }
        );

        function emptyFields() {
            $scope.username = "";
            $scope.password = "";
            $scope.email = "";
        }

        function emptyMessages() {
            $scope.errorMessages = [];
            $scope.successMessages = [];
        }
    }
});