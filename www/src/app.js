var app = angular.module('MaliciousMango', ['ngRoute', 'ngAnimate', 'lr.upload', 'toastr']);

app.run(['tokenRefresher', function(tokenRefresher) {
  tokenRefresher.start(); // will start if neccesary
}]);

app.config(['$routeProvider', '$locationProvider', 
   function ($routeProvider,   $locationProvider) { 
  // $locationProvider.html5Mode(true); // to be used later to remove # in url

  /* This is where all pages is specified */
  $routeProvider 
    .when('/', { 
      templateUrl:'/views/frontpage.html'
    }) 
    .when('/login', { 
      controller:'LoginController',
      templateUrl:'../views/login.html'
    })
    .when('/signup', { 
      controller:'SignupController',
      templateUrl:'../views/signup.html'
    })
    .when('/profile', {
      controller:'ProfileRedirectController',
      templateUrl:'../views/frontpage.html'
    })
    .when('/profile/edit', {
      controller:'ProfileEditController',
      templateUrl:'../views/profileEdit.html'
    })
    .when('/profile/:publicName', {
      controller:'ProfileController',
      templateUrl:'../views/profile.html'
    })
    .otherwise({ 
      redirectTo: '/' 
    }); 
}]);
