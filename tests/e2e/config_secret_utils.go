/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"context"
	"fmt"
	"strings"

	"github.com/onsi/gomega"
	"github.com/vmware/govmomi/object"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	fssh "k8s.io/kubernetes/test/e2e/framework/ssh"
)

// createTestUser util method is used for creating test users
func createTestUser(masterIp string, testUser string, testUserPassword string) error {
	createUser := govcLoginCmd() + "govc sso.user.create -p " + testUserPassword + " " + testUser
	framework.Logf("Create testuser: %s ", createUser)
	result, err := sshExec(sshClientConfig, masterIp, createUser)
	if err != nil && result.Code != 0 {
		fssh.LogResult(result)
		return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
			createUser, masterIp, err)
	}
	return nil
}

// deleteUsersRolesAndPermissions method is used to delete roles and permissions of a test users
func deleteUsersRolesAndPermissions(masterIp string, testUser string, testUserAlias string,
	dataCenter []*object.Datacenter, cluster []string, hosts []string, vms []string, datastores []string) {
	framework.Logf("Delete user permissions")
	deleteUserPermissions(masterIp, testUserAlias, dataCenter, cluster, hosts, vms, datastores)

	framework.Logf("Delete user roles")
	err := deleteUserRoles(masterIp, testUser)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			framework.Logf("No test user roles exist")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}
}

// deleteUserPermissions method is used to delete permissions of a test user
func deleteUserPermissions(masterIp string, testUser string, dataCenter []*object.Datacenter, clusters []string,
	hosts []string, vms []string, datastores []string) {
	err := deleteDataCenterPermissions(masterIp, testUser, dataCenter)
	if err != nil {
		if strings.Contains(err.Error(), "The object or item referred to could not be found") {
			framework.Logf("No datacenter level permissions exist for a testuser")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}

	err = deleteHostsLevelPermission(masterIp, testUser, hosts)
	if err != nil {
		if strings.Contains(err.Error(), "The object or item referred to could not be found") {
			framework.Logf("No host level permissions exist for a testuser")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}

	err = deleteVMsLevelPermission(masterIp, testUser, vms)
	if err != nil {
		if strings.Contains(err.Error(), "The object or item referred to could not be found") {
			framework.Logf("No vm level permissions exist for a testuser")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}

	err = deleteClusterLevelPermission(masterIp, testUser, clusters)
	if err != nil {
		if strings.Contains(err.Error(), "The object or item referred to could not be found") {
			framework.Logf("No cluster level permissions exist for a testuser")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}

	for i := 0; i < len(dataCenter); i++ {
		err = deleteDataStoreLevelPermission(masterIp, testUser, dataCenter[i].InventoryPath, datastores)
		if err != nil {
			if strings.Contains(err.Error(), "The object or item referred to could not be found") {
				framework.Logf("No datastore level permissions exist for a testuser")
			} else {
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
					masterIp, err)
			}
		}
	}
	err = deleteOtherPermissionsFromTestUser(masterIp, testUser)
	if err != nil {
		if strings.Contains(err.Error(), "The object or item referred to could not be found") {
			framework.Logf("No search level permissions exist for a testuser")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}
}

// deleteDataCenterPermissions method is used to delete DataCenter Permissions from a test user
func deleteDataCenterPermissions(masterIp string, testUser string, dataCenter []*object.Datacenter) error {
	for i := 0; i < len(dataCenter); i++ {
		deleteDataCenterPermissions := govcLoginCmd() +
			"govc permissions.remove -principal " + testUser + " " + dataCenter[i].InventoryPath
		framework.Logf("delete datacenter level permissions %s", deleteDataCenterPermissions)
		result, err := sshExec(sshClientConfig, masterIp, deleteDataCenterPermissions)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				deleteDataCenterPermissions, masterIp, err)
		}
	}
	return nil
}

// deleteHostsLevelPermission method is used to delete hosts level permissions from a test user
func deleteHostsLevelPermission(masterIp string, testUser string, hosts []string) error {
	for i := 0; i < len(hosts); i++ {
		deleteHostsLevelPermissions := govcLoginCmd() +
			"govc permissions.remove -principal " + testUser + " " + hosts[i]
		framework.Logf("delete host level permissions %s", deleteHostsLevelPermissions)
		result, err := sshExec(sshClientConfig, masterIp, deleteHostsLevelPermissions)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				deleteHostsLevelPermissions, masterIp, err)
		}
	}
	return nil
}

// deleteVMsLevelPermission method is used to delete vm level permissions from a test user
func deleteVMsLevelPermission(masterIp string, testUser string, vms []string) error {
	for i := 0; i < len(vms); i++ {
		deleteVmsLevelPermissions := govcLoginCmd() + "govc permissions.remove -principal " + testUser +
			" " + vms[i]
		framework.Logf("delete vm level permissions %s", deleteVmsLevelPermissions)
		result, err := sshExec(sshClientConfig, masterIp, deleteVmsLevelPermissions)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				deleteVmsLevelPermissions, masterIp, err)
		}
	}
	return nil
}

// deleteClusterLevelPermission method is used to delete cluster level permissions from a test user
func deleteClusterLevelPermission(masterIp string, testUser string, clusters []string) error {
	for i := 0; i < len(clusters); i++ {
		deleteClusterLevelPermissions := govcLoginCmd() +
			"govc permissions.remove -principal " + testUser + " " + clusters[i]
		framework.Logf("delete cluster level permissions %s", deleteClusterLevelPermissions)
		result, err := sshExec(sshClientConfig, masterIp, deleteClusterLevelPermissions)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				deleteClusterLevelPermissions, masterIp, err)
		}
	}
	return nil
}

// deleteDataStoreLevelPermission method is used to delete datastore level permissions from a test user
func deleteDataStoreLevelPermission(masterIp string, testUser string, dataCenter string, datastores []string) error {
	for i := 0; i < len(datastores); i++ {
		deleteDataStoreLevelPermissions := govcLoginCmd() + "govc permissions.remove -principal " +
			testUser + " '" + datastores[i] + "'"
		framework.Logf("delete datastore level permissions %s", deleteDataStoreLevelPermissions)
		result, err := sshExec(sshClientConfig, masterIp, deleteDataStoreLevelPermissions)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				deleteDataStoreLevelPermissions, masterIp, err)
		}
	}
	return nil
}

// deleteOtherPermissionsFromTestUser method is used to remove permissions from a test user
func deleteOtherPermissionsFromTestUser(masterIp string, testUser string) error {
	deleteOtherPermissions := govcLoginCmd() + "govc permissions.remove -principal " +
		testUser
	framework.Logf("delete other permissions %s", deleteOtherPermissions)
	result, err := sshExec(sshClientConfig, masterIp, deleteOtherPermissions)
	if err != nil && result.Code != 0 {
		fssh.LogResult(result)
		return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
			deleteOtherPermissions, masterIp, err)
	}
	return nil
}

// deleteUserRoles method is used to delete roles of a test user
func deleteUserRoles(masterIp string, testUser string) error {
	roleMap := userRoleMap()
	for key := range roleMap {
		if key != "ReadOnly" {
			deleteRoles := govcLoginCmd() + "govc role.remove " + key + "-" + testUser
			framework.Logf("delete user roles %s", deleteRoles)
			result, err := sshExec(sshClientConfig, masterIp, deleteRoles)
			if err != nil && result.Code != 0 {
				fssh.LogResult(result)
				return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
					deleteRoles, masterIp, err)
			}
		}
	}
	return nil
}

// deleteTestUser method is used to delete config secret test users
func deleteTestUser(masterIp string, testUser string) error {
	deleteUser := govcLoginCmd() + "govc sso.user.rm " + testUser
	framework.Logf("delete test user %s", deleteUser)
	result, err := sshExec(sshClientConfig, masterIp, deleteUser)
	if err != nil && result.Code != 0 {
		fssh.LogResult(result)
		return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
			deleteUser, masterIp, err)
	}
	return nil
}

// createRolesForTestUser method is used to create roles for a test user
func createRolesForTestUser(masterIp string, testUser string) error {
	roleMap := userRoleMap()
	for key, val := range roleMap {
		if key != "ReadOnly" {
			createRoleCmdFortestUser := govcLoginCmd() + "govc role.create " + key + "-" + testUser + " " + val
			framework.Logf("Create roles for test user %s", createRoleCmdFortestUser)
			result, err := sshExec(sshClientConfig, masterIp, createRoleCmdFortestUser)
			if err != nil && result.Code != 0 {
				fssh.LogResult(result)
				return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
					createRoleCmdFortestUser, masterIp, err)
			}
		}
	}
	return nil
}

/*
getDataCenterClusterHostAndVmDetails method is used to fetch data center details, cluster
details, host details and vm details
*/
func getDataCenterClusterHostAndVmDetails(ctx context.Context, masterIp string) ([]*object.Datacenter,
	[]string, []string, []string, []string) {
	// fetch datacenter details
	dataCenters, err := e2eVSphere.getAllDatacenters(ctx)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// fetch cluster details
	clusters, err := getClusterNames(masterIp, dataCenters)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
		masterIp, err)

	// fetch esxi hosts details
	hosts, err := getEsxiHostNames(masterIp, clusters)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
		masterIp, err)

	// fetch vm details
	vms, err := getVmNames(masterIp, dataCenters)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
		masterIp, err)

	// fetch datastore details
	datastores, err := getDatastoreNames(masterIp, dataCenters)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
		masterIp, err)

	return dataCenters, clusters, hosts, vms, datastores
}

// getClusterNames method is used to fetch cluster list
func getClusterNames(masterIp string, dataCenter []*object.Datacenter) ([]string, error) {
	var clusDetails, clusterList, clusterNames []string
	framework.Logf("Fetching cluster details")
	for i := 0; i < len(dataCenter); i++ {
		clusterFolder := govcLoginCmd() + "govc ls " + dataCenter[i].InventoryPath
		clusterFolderNameResult, err := sshExec(sshClientConfig, masterIp, clusterFolder)
		if err != nil && clusterFolderNameResult.Code != 0 {
			fssh.LogResult(clusterFolderNameResult)
			return nil, fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				clusterFolder, masterIp, err)
		}
		if clusterFolderNameResult.Stdout != "" {
			clusDetails = strings.Split(clusterFolderNameResult.Stdout, "\n")
		}
		clusterPathName := ""
		for i := 0; i < len(clusDetails)-1; i++ {
			if strings.Contains(clusDetails[i], "host") {
				clusterPathName = clusDetails[i]
				break
			}
		}
		clusterGroup := govcLoginCmd() + "govc ls " + clusterPathName
		clusterGroupResult, err := sshExec(sshClientConfig, masterIp, clusterGroup)
		if err != nil && clusterGroupResult.Code != 0 {
			fssh.LogResult(clusterGroupResult)
			return nil, fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				clusterGroup, masterIp, err)
		}
		if clusterGroupResult.Stdout != "" {
			clusterNames = strings.Split(clusterGroupResult.Stdout, "\n")
		}
		cluster := govcLoginCmd() + "govc ls " + clusterNames[0] + " | sort"
		clusterResult, err := sshExec(sshClientConfig, masterIp, cluster)
		if err != nil && clusterResult.Code != 0 {
			fssh.LogResult(clusterResult)
			return nil, fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				cluster, masterIp, err)
		}
		clusDetails = nil
		if !strings.Contains(clusterResult.Stdout, "10.") {
			if clusterResult.Stdout != "" {
				clusListTemp := strings.Split(clusterResult.Stdout, "\n")
				clusDetails = append(clusDetails, clusListTemp...)
			}
			for i := 0; i < len(clusDetails)-1; i++ {
				clusterList = append(clusterList, clusDetails[i])
			}
			clusDetails = nil
		} else {
			for i := 0; i < len(clusterNames)-1; i++ {
				clusterList = append(clusterList, clusterNames[i])
			}
			clusDetails = nil
		}
	}
	return clusterList, nil
}

// getEsxiHostNames method is used to fetch esxi hosts details
func getEsxiHostNames(masterIp string, cluster []string) ([]string, error) {
	var hostsList, hostList []string
	framework.Logf("Fetching ESXi host details")
	for i := 0; i < len(cluster); i++ {
		hosts := govcLoginCmd() + "govc ls " + cluster[i] + " " + " | grep 10."
		hostsResult, err := sshExec(sshClientConfig, masterIp, hosts)
		if err != nil && hostsResult.Code != 0 {
			fssh.LogResult(hostsResult)
			return nil, fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				hosts, masterIp, err)
		}
		if hostsResult.Stdout != "" {
			hostListTemp := strings.Split(hostsResult.Stdout, "\n")
			hostList = append(hostList, hostListTemp...)
		}
		for i := 0; i < len(hostList)-1; i++ {
			hostsList = append(hostsList, hostList[i])
		}
		hostList = nil
	}
	return hostsList, nil
}

// getVmNames method is used to fetch vm details
func getVmNames(masterIp string, dataCenter []*object.Datacenter) ([]string, error) {
	var vmsList, vmList []string
	framework.Logf("Fetching VM details")
	for i := 0; i < len(dataCenter); i++ {
		vms := govcLoginCmd() + "govc ls " + dataCenter[i].InventoryPath + "/vm" + " " + "| grep 'k8s\\|haproxy'"
		vMsResult, err := sshExec(sshClientConfig, masterIp, vms)
		if err != nil && vMsResult.Code != 0 {
			fssh.LogResult(vMsResult)
			return nil, fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				vms, masterIp, err)
		}
		if vMsResult.Stdout != "" {
			vmListTemp := strings.Split(vMsResult.Stdout, "\n")
			vmList = append(vmList, vmListTemp...)
		}
		for i := 0; i < len(vmList)-1; i++ {
			vmsList = append(vmsList, vmList[i])
		}
		vmList = nil
	}
	return vmsList, nil
}

// getDatastoreNames method is used to fetch datastore details
func getDatastoreNames(masterIp string, dataCenter []*object.Datacenter) ([]string, error) {
	var dsList, datastores []string
	framework.Logf("Fetching datastore details")
	for i := 0; i < len(dataCenter); i++ {
		ds := govcLoginCmd() + "govc ls " + dataCenter[i].InventoryPath + "/datastore"
		dsResult, err := sshExec(sshClientConfig, masterIp, ds)
		if err != nil && dsResult.Code != 0 {
			fssh.LogResult(dsResult)
			return nil, fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				ds, masterIp, err)
		}
		if dsResult.Stdout != "" {
			dsListTemp := strings.Split(dsResult.Stdout, "\n")
			dsList = append(dsList, dsListTemp...)
		}
		for i := 0; i < len(dsList)-1; i++ {
			datastores = append(datastores, dsList[i])
		}
		dsList = nil
	}
	return datastores, nil
}

// setDataCenterLevelPermission is used to set data center level permissions for test user
func setDataCenterLevelPermission(masterIp string, dataCenter string, testUser string,
	propagateVal string, readOnlyRole string) error {
	setPermissionForDataCenter := govcLoginCmd() + "govc permissions.set -principal " + testUser +
		" -propagate=" + propagateVal + " -role " + readOnlyRole + " " + dataCenter + " | tr -d '\n'"
	result, err := sshExec(sshClientConfig, masterIp, setPermissionForDataCenter)
	if err != nil && result.Code != 0 {
		fssh.LogResult(result)
		return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
			setPermissionForDataCenter, masterIp, err)
	}

	return nil
}

// setHostLevelPermission is used to set host level permissions for test user
func setHostLevelPermission(masterIp string, testUser string, hosts []string, propagateVal string,
	readOnlyRole string) error {
	for i := 0; i < len(hosts); i++ {
		setPermissionForHosts := govcLoginCmd() + "govc permissions.set -principal " +
			testUser + " -propagate=" + propagateVal + " -role " + readOnlyRole + " " + hosts[i] + "| tr -d '\n'"
		result, err := sshExec(sshClientConfig, masterIp, setPermissionForHosts)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				setPermissionForHosts, masterIp, err)
		}
	}

	return nil
}

// setVMLevelPermission is used to set vm level permissions for test user
func setVMLevelPermission(masterIp string, testUserAlias string, testUser string, vms []string,
	propagateVal string, vmRole string) error {
	for i := 0; i < len(vms); i++ {
		setPermissionForK8sVms := govcLoginCmd() + "govc permissions.set -principal " +
			testUserAlias + " -propagate=" + propagateVal + " -role " + vmRole + "-" + testUser +
			" " + vms[i] + " | tr -d '\n'"
		result, err := sshExec(sshClientConfig, masterIp, setPermissionForK8sVms)
		framework.Logf("Vm level permissions %s", setPermissionForK8sVms)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				setPermissionForK8sVms, masterIp, err)
		}
	}
	return nil
}

// setClusterLevelPermission is used to set cluster level permissions for test user
func setClusterLevelPermission(masterIp string, testUserAlias string, testUser string, cluster string,
	propagateVal string, hostRole string) error {
	setPermissionForCluster := govcLoginCmd() + "govc permissions.set -principal " +
		testUserAlias + " -propagate=" + propagateVal + " -role " + hostRole + "-" +
		testUser + " " + cluster + " | tr -d '\n'"
	framework.Logf("Cluster level permissions %s", setPermissionForCluster)
	result, err := sshExec(sshClientConfig, masterIp, setPermissionForCluster)
	if err != nil && result.Code != 0 {
		fssh.LogResult(result)
		return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
			setPermissionForCluster, masterIp, err)
	}
	return nil
}

// setDataStoreLevelPermission is used to set datastore level permissions for test user
func setDataStoreLevelPermission(masterIp string, testUserAlias string, testUser string, dataCenter string,
	datastores []string, propagateVal string, datastoreRole string) error {
	for i := 0; i < len(datastores); i++ {
		setPermissionForDataStore := govcLoginCmd() + "govc permissions.set -principal " +
			testUserAlias + " " + "-propagate=" + propagateVal + " -role " + datastoreRole + "-" + testUser + " '" +
			datastores[i] + "'"
		framework.Logf("Datastore level permissions %s", setPermissionForDataStore)
		result, err := sshExec(sshClientConfig, masterIp, setPermissionForDataStore)
		if err != nil && result.Code != 0 {
			fssh.LogResult(result)
			return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
				setPermissionForDataStore, masterIp, err)
		}
	}
	return nil
}

// setSearchlevelPermission method is used to set search level permissions
func setSearchlevelPermission(masterIp string, testUserAlias string, testUser string,
	propagateVal string, searchRole string) error {
	setSearchLevelPermission := govcLoginCmd() + "govc permissions.set -principal " + testUserAlias +
		" -propagate=" + propagateVal + " -role " + searchRole + "-" + testUser + " /"
	framework.Logf("Search level permissions %s", setSearchLevelPermission)
	result, err := sshExec(sshClientConfig, masterIp, setSearchLevelPermission)
	if err != nil && result.Code != 0 {
		fssh.LogResult(result)
		return fmt.Errorf("couldn't execute command: %s on host: %v , error: %s",
			setSearchLevelPermission, masterIp, err)
	}
	return nil
}

// createCsiVsphereSecret method is used to create csi vsphere secret file
func createCsiVsphereSecret(client clientset.Interface, ctx context.Context, testUser string,
	password string, csiNamespace string) {
	currentSecret, err := client.CoreV1().Secrets(csiNamespace).Get(ctx, configSecret, metav1.GetOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	originalConf := string(currentSecret.Data[vSphereCSIConf])
	vsphereCfg, err := readConfigFromSecretString(originalConf)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	vsphereCfg.Global.User = testUser
	vsphereCfg.Global.Password = password
	modifiedConf, err := writeConfigToSecretString(vsphereCfg)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	framework.Logf("Updating the secret to reflect new conf credentials")
	currentSecret.Data[vSphereCSIConf] = []byte(modifiedConf)
	_, err = client.CoreV1().Secrets(csiNamespace).Update(ctx, currentSecret, metav1.UpdateOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

/*
createTestUserAndAssignRolesPrivileges method is used to create test user, assign
roles and privileges to test user
*/
func createTestUserAndAssignRolesPrivileges(masterIp string, configSecretTestUser string,
	configSecretTestUserPassword string, configSecretTestUserAlias string, propagateVal string,
	dataCenters []*object.Datacenter, clusters []string, hosts []string,
	vms []string, datastores []string) {
	roleMap := userRoleMap()

	framework.Logf("Create TestUser")
	err := createTestUser(masterIp, configSecretTestUser, configSecretTestUserPassword)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
		masterIp, err)

	framework.Logf("Create roles for TestUser")
	err = createRolesForTestUser(masterIp, configSecretTestUser)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
		masterIp, err)

	for key := range roleMap {
		if strings.Contains(key, "VM") {
			framework.Logf("Assign vm level permissions")
			err = setVMLevelPermission(masterIp, configSecretTestUserAlias, configSecretTestUser, vms,
				propagateVal, key)
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
		if strings.Contains(key, "HOST") {
			framework.Logf("Assign cluster level permissions")
			for i := 0; i < len(clusters); i++ {
				err = setClusterLevelPermission(masterIp, configSecretTestUserAlias, configSecretTestUser,
					clusters[i], propagateVal, key)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
					masterIp, err)
			}
		}
		if strings.Contains(key, "DATASTORE") {
			framework.Logf("Assign datastores level permissions")
			for i := 0; i < len(dataCenters); i++ {
				err = setDataStoreLevelPermission(masterIp, configSecretTestUserAlias, configSecretTestUser,
					dataCenters[i].InventoryPath, datastores, propagateVal, key)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
					masterIp, err)
			}
		}
		if strings.Contains(key, "SEARCH") {
			framework.Logf("Assign search level permissions")
			err = setSearchlevelPermission(masterIp, configSecretTestUserAlias, configSecretTestUser,
				propagateVal, key)
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
		if strings.Contains(key, "ReadOnly") {
			framework.Logf("Assign datacenter level read-only permissions")
			for i := 0; i < len(dataCenters); i++ {
				err = setDataCenterLevelPermission(masterIp, dataCenters[i].InventoryPath, configSecretTestUserAlias,
					propagateVal, key)
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
					masterIp, err)
			}

			framework.Logf("Assign host level read-only permissions")
			err = setHostLevelPermission(masterIp, configSecretTestUserAlias, hosts, propagateVal, key)
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}
}

/*
deleteTestUserAndRemoveRolesPrivileges method is used to delete test user and to remove assigned
roles and privileges to test user
*/
func deleteTestUserAndRemoveRolesPrivileges(masterIp string, configSecretTestUser string,
	configSecretTestUserPassword string, configSecretTestUserAlias string, propagateVal string,
	dataCenters []*object.Datacenter, clusters []string, hosts []string,
	vms []string, datastores []string) {
	framework.Logf("Delete users roles and permissions")
	deleteUsersRolesAndPermissions(masterIp, configSecretTestUser, configSecretTestUserAlias, dataCenters,
		clusters, hosts, vms, datastores)

	framework.Logf("Delete Test user")
	err := deleteTestUser(masterIp, configSecretTestUser)
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			framework.Logf("test user doesn't exist")
		} else {
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "couldn't execute command on host: %v , error: %s",
				masterIp, err)
		}
	}
}

// userRoleMap util method returns a map of roles required for a test user
func userRoleMap() map[string]string {
	roleMap := make(map[string]string)
	roleMap["CNS-DATASTORE"] = "Datastore.FileManagement"
	roleMap["CNS-HOST-CONFIG-STORAGE"] = "Host.Config.Storage"
	roleMap["CNS-VM"] = "VirtualMachine.Config.AddExistingDisk VirtualMachine.Config.AddRemoveDevice"
	roleMap["CNS-SEARCH-AND-SPBM"] = "Cns.Searchable StorageProfile.View"
	roleMap["ReadOnly"] = "ReadOnly"

	return roleMap
}
