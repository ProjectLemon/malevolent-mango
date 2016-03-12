var app = angular.module('MaliciousMango', ['ngRoute', 'ngAnimate']);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) { 
  // $locationProvider.html5Mode(true); to be used later to remove # in url

  /* This is where all pages is specified */
  $routeProvider 
    .when('/', { 
      templateUrl:'/views/frontpage.html'
    }) 
    .when('/login', { 
      controller:'LoginController',
      templateUrl:'/views/login.html'
    }) 
    .otherwise({ 
      redirectTo: '/' 
    }); 
    
  
}]);
