var app = angular.module("MariaDBApp", ['ngRoute', 'ngFileUpload', 'ui.bootstrap'])
    .config(function ($routeProvider,$locationProvider) {


        $locationProvider.hashPrefix(''); // TEMP


        $routeProvider
            .when('/', {
                templateUrl: "view/HomeView.html",
                controller: "HomeController",
                activetab: 'home'
            })
            .when("/login", {
                templateUrl: "view/LoginView.html",
                controller: "LoginController",
                activetab: 'login'
            })
            .when("/register", {
                templateUrl: "view/RegisterView.html",
                controller: "RegisterController",
                activetab: 'register'
            })
            .when('/profile', {
                templateUrl: "view/ProfileView.html",
                controller: "ProfileController",
                activetab: 'profile'
            })
            .when('/upload', {
                templateUrl: "view/UploadView.html",
                controller: "UploadController",
                activetab: 'upload'
            })
            .when("/logout", {
                templateUrl: "view/LogoutView.html",
                controller: "LogoutController",
                activetab: 'logout'
            })
            .when("/toprated", {
                templateUrl: "view/TopRatedView.html",
                controller: "TopRatedController",
                activetab: 'toprated'
            })
            .when("/hot", {
                templateUrl: "view/HotView.html",
                controller: "HotController",
                activetab: 'hot'
            })
            .when("/user", {
                templateUrl: "view/UserView.html",
                controller: "UserController",
                activetab: 'user'
            })
            .otherwise({
                redirectTo: "/"
            });
    });