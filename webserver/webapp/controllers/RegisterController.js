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



                if (data.user) {
                    //console.info(data.user);
                    var user = data.user;
                    if (user.id && user.email && user.username) {
                        LocalStorage.setToken(data.token);
                        LocalStorage.setUser(data.user);

                        $window.location.href = '#/';
                        return;
                    }
                }
                $scope.errorMessages.push("Something went wrong. Please try again.");

            }, function (response) {
                console.log(response);
                if (response && response.data) {
                    var data = response.data;

                    if (data && data.message) {
                        $scope.errorMessages.push(data.message);
                        return;
                    }
                }
                $scope.errorMessages.push("An error occurred. Please try again.");
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