app.controller('ProfileRedirectController', ['$window', function ($window) {
  if ($window.sessionStorage.token) {
    $window.location.href = '#/edit';
  } else {
    $window.location.href = '/';
  }
}]);