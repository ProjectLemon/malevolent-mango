/**
 * ProfileEditController handles the edit page for the user profile.
 * Will get info from server, handle edits, and then push changes back to server
 */
app.controller('ProfileEditController', ['$scope', '$http', '$window', '$timeout', function ($scope, $http, $window, $timeout) {
  // Declare variables
  $scope.user = {
      FullName: 'Nathan Drake',
      Email: 'testing@example.com',
      Phone: '073-902301',
      Description: 'Some description about something or other. Oh look at me, I\'m just writing enough text to get a new line. Lorem ipsum dolar cofal(?) and all that shit.',
      
      ProfileIcon: 'img/testFace.png',
      ProfileHeader: 'img/testBG.png',
      
      Pdfs: [
          {title: 'Portfolio', path: 'pdfs/portfolio1.pdf'}, 
          {title: 'Resumé', path: 'pdfs/resume1.pdf'},
          {title: '+'}
      ]
  };
  $scope.message = ''; 
  $scope.loading = {header: false, icon: false};

  
  //Do a http request to server
  $http.get('/api/profile').then(
    //If success
    // Get user information from server and puts it in the user variable
    function (response) {
      //$scope.user = response.data;
    },
    //If Error
    // Display message that the user is not found
    function(response) {
      $scope.message = 'User is not found'; 
    }
  )
  
  $scope.changeBackground = function(response) {
    $scope.user.ProfileHeader = response.data;
  }
  $scope.changeProfileIcon = function(response) {
    $scope.user.ProfileIcon = response.data;
  }
  $scope.startUploadHeader = function(response) {
    $scope.loading.header = true;
  }
  $scope.startUploadIcon = function(response) {
    $scope.loading.icon = true;
  }
  $scope.endUploadHeader = function(response) {
    $scope.loading.header = false;
  }
  $scope.endUploadIcon = function(response) {
    $scope.loading.icon = false;
  }
}]);