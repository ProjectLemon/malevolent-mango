/**
 * pdfTab
 */
app.directive('pdfTab', function() {
  return {
    restrict: 'EA',
    scope: {
      title: '@'
    },
    
    link: function(scope, element, attrs) {
    }
  };
});