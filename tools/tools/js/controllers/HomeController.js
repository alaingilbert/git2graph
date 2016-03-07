var app = angular.module('app', ['ngSanitize', 'ui.bootstrap', 'LocalStorageModule',]);

app.controller('HomeController',
  function($scope, $uibModal, localStorageService)
  {

    var recreateIds = function() {
      _.map($scope.tree, function(node, idx) {
        node.id = idx.toString();
      });
    };

    $scope.btnDeleteNodeClicked = function() {
      $scope.tree.splice($scope.selectedNode.id, 1);
      recreateIds();
      $scope.selectedNode = null;
    };


    $scope.btnMoveDownClicked = function(path, $index) {
      var point = path.splice($index, 1)[0];
      path.splice($index + 1, 0, point);
    };


    $scope.btnMoveUpClicked = function(path, $index) {
      var point = path.splice($index, 1)[0];
      path.splice($index - 1, 0, point);
    };


    $scope.btnRemovePointClicked = function(path, $index) {
      path.splice($index, 1);
    };


    $scope.btnAddPointClicked = function(path) {
      var point = [0, 0, 0];
      path.push(point);
    };


    $scope.btnSaveClicked = function() {
      var modalInstance = $uibModal.open({
        animation: true,
        template:
          '<div class="modal-content">' +
          '  <div class="modal-header">' +
          '    <button type="button" class="close" ng-click="cancel()">x</button>' +
          '    <h4 class="modal-title" id="mySmallModalLabel">Save</h4>' +
          '  </div>' +
          '  <div class="modal-body">' +
          '    <input ng-model="input.name" type="text" class="form-control" ng-enter="btnSaveClicked()" focus-me required />' +
          '  </div>' +
          '  <div class="modal-footer">' +
          '    <input type="button" value="Save" class="btn btn-primary" ng-click="btnSaveClicked()" />' +
          '    <input type="button" value="Close" class="btn btn-default" ng-click="cancel()" />' +
          '  </div>' +
          '</div>',
        controller: function($scope, tree, $uibModalInstance, projectName) {
          $scope.input = {};
          $scope.input.name = projectName;
          $scope.btnSaveClicked = function() {
            if (_.isEmpty($scope.input.name)) {
              return;
            }
            var names = localStorageService.get('names') || {};
            names[$scope.input.name] = tree;
            localStorageService.set('names', names);
            $uibModalInstance.close();
          };
          $scope.cancel = function() {
            $uibModalInstance.dismiss('cancel');
          };
        },
        size: 'sm',
        resolve: {
          projectName: function() { return $scope.input.projectName; },
          tree: function() { return $scope.tree; }
        },
      });
    };

    $scope.btnLoadClicked = function() {
      var modalInstance = $uibModal.open({
        animation: true,
        template:
          '<div class="modal-content">' +
          '  <div class="modal-header">' +
          '    <button type="button" class="close" ng-click="cancel()">x</button>' +
          '    <h4 class="modal-title" id="mySmallModalLabel">Load</h4>' +
          '  </div>' +
          '  <div class="modal-body">' +
          '    <select ng-model="input.name" ng-options="name for name in names" class="form-control" ng-enter="btnLoadClicked()" focus-me></select>' +
          '  </div>' +
          '  <div class="modal-footer">' +
          '    <input type="button" value="Load" class="btn btn-primary" ng-click="btnLoadClicked()" />' +
          '    <input type="button" value="Close" class="btn btn-default" ng-click="cancel()" />' +
          '  </div>' +
          '</div>',
        controller: 'LoadModalController',
        size: 'sm',
      }).result.then(function(res) {
        if (res) {
          $scope.input.projectName = res.projectName;
          $scope.tree = res.data;
        }
      });
    };


    var highlightJson = function(json) {
      json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
      return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
        var cls = 'number';
        if (/^"/.test(match)) {
          if (/:$/.test(match)) {
            cls = 'key';
          } else {
            cls = 'string';
          }
        } else if (/true|false/.test(match)) {
          cls = 'boolean';
        } else if (/null/.test(match)) {
          cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
      });
    };

    $scope.highlight = function(obj) {
      return highlightJson(angular.toJson(obj, true))
    };


    (function constructor() {
      $scope.selectedNode = null;
      $scope.inputFile = '';
      $scope.testFile = '';
      $scope.shellFile = '';
      $scope.tree = [];
      $scope.input = {};

      $scope.colors = [
        '#5aa1be',
        '#c065b8',
        '#c0ab5f',
        '#59bc95',
        '#7a63be',
        '#c0615b',
        '#73bb5e',
        '#6ee585',
        '#7088e8',
        '#eb77a3',
        '#c2e675',
      ];

    })();

  }
);
