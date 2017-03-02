describe('Testing the app', function () {

    // Before each
    var mockWindow, apiService, httpBackend;
    beforeEach(function () {
        module('MariaDBApp');
    });

    beforeEach(inject(function ($window, $httpBackend, ApiService) {
        mockWindow = $window;
        httpBackend = $httpBackend;
        apiService = ApiService;
    }));


    // Tests

    it('should login using ApiService and return 200 status', inject(function ($http) {

        var $scope = {};

        /* Code Under Test */
        apiService.login("user", "pass").then(
            function (response) {
                console.info(response);

                if (response) {
                    if (response.token && response.user && response.user.id) {
                        $scope.valid = true;
                    }
                }
                $scope.response = response;
            }, function (error) {
                console.error(error);
                $scope.valid = false;

            }
        );
        /* End */

        httpBackend
            .when('POST', './token-auth')
            .respond(200, {
                'token': 'HFAIF&(SD&F(DS&F(*DSF&SDF^%SDF&^SD^F$&SD$F^S^D%F$&SD%F',
                'expires_on': 43095835,
                'user': {
                    'id': 1,
                    'username': 'user',
                    'email': 'user@test.com',
                    'createdAt': 9879456917
                }
            });

        httpBackend.flush();

        expect($scope.valid).toBe(true);
        var object = {
            'token': 'HFAIF&(SD&F(DS&F(*DSF&SDF^%SDF&^SD^F$&SD$F^S^D%F$&SD%F',
            'expires_on': 43095835,
            'user': {
                'id': 1,
                'username': 'user',
                'email': 'user@test.com',
                'createdAt': 9879456917
            }
        };
        expect($scope.response).toEqual(object);

    }));

});