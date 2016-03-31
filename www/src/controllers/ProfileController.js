app.controller('ProfileController', ['$scope', '$http', '$window', function ($scope, $http, $window) {
  // Declare variables
  $scope.user = {
      FullName: 'Nathan Drake',
      Email: 'testing@example.com',
      Phone: '073-902301',
      Description: 'Somtha asfkjnsadfjkdfsajkklöjasdf',
      
      ProfileIcon: 'images/test.png',
      ProfileHeader: 'images/testBG.png',
  };
  $scope.message = ''; 

  
  //Do a http request to server
  $http.get('/api/profile').then(
    //If success
    // Get user information from server and puts it in the user variable
    function (response) {
      $scope.user = response.data;
    },
    //If Error
    // Display message that the user is not found
    function(response) {
      $scope.message = 'User is not found'; 
    }
  )
}]);