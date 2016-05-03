/**
 * pdfViewer shows a pdf and update it when source change
 * Should have two attributes, pdfs and current. Current is only
 * a number but pdfs should be an array of pdfs in form of
 * [{Title: 'title', Path: 'som/path'}, ...]
 */
app.directive('pdfViewer', ['$compile', function($compile) {
  return {
    restrict: 'E',
    scope: {
      pdfs: '=',
      current: '@'
    },

    link: function(scope, element, attrs) {

      var currentElement = element;
      scope.$parent.$watch(scope.current, function(newValue, oldValue) {
        var html;
        if (scope.pdfs && scope.pdfs.length > 0 && scope.pdfs[newValue] != null) {
          
          var pdf = scope.pdfs[newValue];
          html = '<embed src="'+pdf.Path+'" type="application/pdf">'

        } else {
          html = '<p>No pdf</p>'
        }
        var e = $compile(html)(scope);
        currentElement.replaceWith(e);
        currentElement = e;
      });
    }
  };
}]);
