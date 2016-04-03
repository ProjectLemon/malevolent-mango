/**
 * Redirects users trying to navigate to a unused url: #/profile
 */
app.controller('ProfileRedirectController', ['$window', function ($window) {
  if ($window.sessionStorage.token) {
    $window.location.href = '#/profile/edit';
  } else {
    $window.location.href = '/';
  }
}]);