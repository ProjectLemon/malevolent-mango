/**
 * ProfileEditController handles the edit page for the user profile.
 * Will get info from server, handle edits, and then push changes back to server
 */
app.controller('ProfileEditController', ['$scope', '$http', '$window', '$timeout', 
                                function ($scope,   $http,   $window,   $timeout) {
  // Declare variables
  $scope.user = { // Placeholder
      FullName: 'Full Name',
      EMail: 'Your email adress',
      Phone: 'Your phone number',
      Description: 'Input your description here',

      ProfileIcon: 'img/profileDefault.png',
      ProfileHeader: 'img/backgroundDefault.png',

      PDFs: [
      ]
  };
  $scope.message = '';
  $scope.loading = {header: false, icon: false};


  //Do a http request to server
  $http.get('/api/profile/get-edit').then(
    //If success
    // Get user information from server and puts it in the user variable
    function (response) {
      if (response.status == 204) {
        
      } else {
        if (response.data != '') {
          $scope.user = response.data;
        }
      }
    },
    //If Error
    function(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content'; 
      } else {
        $scope.message = 'User is not found'; 
      }
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
  
  $scope.publish = function() {
    $http.post('api/profile/save', $scope.user).then(
      function success(response) {
        $scope.message = 'Saved Success'; 
      },
      function error(response) {
        $scope.message = 'Saved failed';
      }
    );
  }
}]);
