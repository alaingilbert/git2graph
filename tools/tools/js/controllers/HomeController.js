var app = angular.module('app', ['ngSanitize', 'ui.bootstrap', 'LocalStorageModule',]);

app.controller('HomeController',
  function($scope, $uibModal, localStorageService)
  {

    var recreateIds = function() {
      _.map($scope.tree, function(node, idx) {
        node.id = idx.toString();
      });
    };

    var removeFromChildrenParents = function(parentId) {
      _.each($scope.tree, function(node) {
        delete node.parentsPaths[parentId];
      });
      var nbNodes = $scope.tree.length;
      _.each($scope.tree, function(node) {
        _.each(node.parentsPaths, function(path, key) {
          if (key > parentId) {
            node.parentsPaths[key-1] = path;
          }
        });
        delete node.parentsPaths[nbNodes];
      });
    };

    $scope.btnDeleteNodeClicked = function() {
      var idToRemove = $scope.selectedNode.id;
      $scope.tree.splice(idToRemove, 1);
      recreateIds();
      removeFromChildrenParents(idToRemove);
      _.each($scope.tree, function(node) {
        _.each(node.parents, function(nodeParent, idx) {
          if (nodeParent == idToRemove) {
            node.parents.splice(idx, 1);
          } else if (nodeParent > idToRemove) {
            node.parents[idx]--;
          }
        });
        _.each(node.parentsPaths, function(path, parentId) {
          _.each(path.path, function(point, idx) {
            if (point[1] > idToRemove) {
              node.parentsPaths[parentId].path[idx][1]--;
            }
          });
        });
      });
      $scope.selectedNode = null;
    };

    $scope.btnDeletePathClicked = function(key) {
      delete $scope.selectedNode.parentsPaths[key];
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


    $scope.toJson = function(tree) {
      return angular.toJson(tree, true);
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
