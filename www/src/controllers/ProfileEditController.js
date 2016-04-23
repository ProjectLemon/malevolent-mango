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
  $scope.currentPDF = 0;
  $scope.loading = {header: false, icon: false, pdf: false};
  $scope.maxLength = {
    fullName: 70,
    email: 50,
    phone: 80,
    description: 150,
    profileIcon: 150,
    profileHeader: 360,
    pdfs: 21844
  }


  //Do a http request to server
  $http.get('/api/profile/get-edit').then(
    //If success
    // Get user information from server and puts it in the user variable
    function (response) {
      if (response.status == 204) {

      } else {
        if (response.data != '') {
          $scope.user = response.data;
          $scope.currentPDF = 0;
        }
      }
    },
    //If Error
    function(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content';

      } else if (response.status == 400) {
        $scope.message = 'You are not logged in';
        $window.location.href = '/'; // return to start page

      } else if (response.status == 413) {

      } else {
        $scope.message = 'User is not found';
      }
    }
  );

  $scope.publish = function() {
    $http.post('api/profile/save', $scope.user).then(
      function success(response) {
        $scope.message = 'Saved Success';
      },
      function error(response) {
        $scope.message = 'Saved failed';
      }
    );
  };

  $scope.changeBackground = function(response) {
    $scope.user.ProfileHeader = response.data;
  }
  $scope.changeProfileIcon = function(response) {
    $scope.user.ProfileIcon = response.data //This is never run?
  }
  $scope.changeProfileIcon = function(response) {
    try {
        $scope.user.PDFs.push({Title: 'Unnamed', Path: response.data});
        $scope.currentPDF = $scope.user.PDFs.length-1;
    } catch (e) {}
  }

  $scope.startUploadHeader = function(response) {
    $scope.loading.header = true;
  }
  $scope.startUploadIcon = function(response) {
    $scope.loading.icon = true;
  }
  $scope.startUploadPdf = function(response) {
    $scope.loading.pdf = true;
  }
  $scope.endUploadHeader = function(response) {
    $scope.loading.header = false;
  }
  $scope.endUploadIcon = function(response) {
    $scope.loading.icon = false;
  }
  $scope.endUploadPDF = function(response) {
    $scope.loading.pdf = false;
  }
}]);
