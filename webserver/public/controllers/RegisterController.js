app.controller('RegisterController', function ($scope, ApiService) {
    $scope.errorMessages = [];
    $scope.successMessages = [];

    $scope.register = function () {
        emptyMessages();

        ApiService.register($scope.username, $scope.email, $scope.password).then(
            function (data) {

                $scope.successMessages.push("Registration successful");
                emptyFields();

            }, function (response) {
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