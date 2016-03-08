var app = angular.module('MaliciousMango', ['ngRoute']);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) { 
  // $locationProvider.html5Mode(true); to be used later to remove # in url

  $routeProvider 
    .when('/login', { 
      controller:'LoginController',
      templateUrl:'/views/login.html'
    }) 
    .otherwise({ 
      redirectTo: '/' 
    }); 
    
  
}]);