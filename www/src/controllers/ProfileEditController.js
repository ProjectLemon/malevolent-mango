/**
 * ProfileEditController handles the edit page for the user profile.
 * Will get info from server, handle edits, and then push changes back to server
 */
app.controller('ProfileEditController', ['$scope', '$http', '$window', '$location', '$timeout', '$interval',
                                function ($scope,   $http,   $window,   $location,   $timeout,   $interval) {
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
  $scope.saved = true;
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
  

  /* Notify user that changes not saved (on window close) */
  var checkIfSaved = function() {
    $scope.$watch('user', function(newValue, oldValue) {
        if (newValue != oldValue) {
          $scope.saved = false;
        }
      }, true); // true to check for value equality
    
    var preventClosingIfNotSaved = function(event) {
      if ($scope.saved == false) {
        event.returnValue = "Warning. You have not published your changes. Leaving this page will remove all changes."
      }
    };
    if (window.addEventListener) {
      window.addEventListener("beforeunload", preventClosingIfNotSaved);
    } else {
      // For IE browsers
      window.attachEvent("onbeforeunload", preventClosingIfNotSaved);
    }
  };



  /* Do a http request to server*/
  $http.get('/api/profile/get-edit').then(

    // Get user information from server and puts it in the user variable
    function success(response) {
      if (response.data != '') {
        $scope.user = response.data;
        $scope.currentPDF = 0;
        checkIfSaved();
      }
    },
    function error(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content';

      } else if (response.status == 400) {
        $scope.message = 'You are not logged in';
        $location.path('/'); // return to start page

      } else if (response.status == 413) {

      } else {
        $scope.message = 'User is not found';
      }
    }
  );

  /* Publish changes to server */
  $scope.publish = function() {
    $http.post('api/profile/save', $scope.user).then(
      function success(response) {
        var oldMessage = $scope.message;
        $scope.message = 'Saved Success';
        $scope.saved = true;
        $interval(function() {$scope.message = oldMessage;}, 5*1000); // 5 sec
      },
      function error(response) {
        var oldMessage = $scope.message;
        $scope.message = 'Saved failed';
        $interval(function() {$scope.message = oldMessage;}, 5*1000); // 5 sec
      }
    );
  };  

  $scope.changeBackground = function(response) {
    $scope.user.ProfileHeader = response.data;
  }
  $scope.changeProfileIcon = function(response) {
    $scope.user.ProfileIcon = response.data //This is never run?
  }
  $scope.addPDF = function(response) {
    if ($scope.user.PDFs == null) {
      $scope.user.PDFs = [];
    }
    $scope.user.PDFs.push({Title: 'Unnamed', Path: response.data});
    $scope.currentPDF = $scope.user.PDFs.length-1;
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
