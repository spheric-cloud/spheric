<p>Packages:</p>
<ul>
<li>
<a href="#core.spheric.cloud%2fv1alpha1">core.spheric.cloud/v1alpha1</a>
</li>
</ul>
<h2 id="core.spheric.cloud/v1alpha1">core.spheric.cloud/v1alpha1</h2>
<div>
<p>Package v1alpha1 is the v1alpha1 version of the API.</p>
</div>
Resource Types:
<ul><li>
<a href="#core.spheric.cloud/v1alpha1.Disk">Disk</a>
</li><li>
<a href="#core.spheric.cloud/v1alpha1.DiskType">DiskType</a>
</li><li>
<a href="#core.spheric.cloud/v1alpha1.Fleet">Fleet</a>
</li><li>
<a href="#core.spheric.cloud/v1alpha1.Instance">Instance</a>
</li><li>
<a href="#core.spheric.cloud/v1alpha1.InstanceType">InstanceType</a>
</li><li>
<a href="#core.spheric.cloud/v1alpha1.Network">Network</a>
</li><li>
<a href="#core.spheric.cloud/v1alpha1.Subnet">Subnet</a>
</li></ul>
<h3 id="core.spheric.cloud/v1alpha1.Disk">Disk
</h3>
<div>
<p>Disk is the Schema for the disks API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Disk</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DiskSpec">
DiskSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>typeRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>TypeRef references the DiskClass of the Disk.</p>
</td>
</tr>
<tr>
<td>
<code>instanceRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalUIDReference">
LocalUIDReference
</a>
</em>
</td>
<td>
<p>InstanceRef references the using instance of the Disk.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.ResourceList">
ResourceList
</a>
</em>
</td>
<td>
<p>Resources is a description of the Disk&rsquo;s resources and capacity.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DiskStatus">
DiskStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DiskType">DiskType
</h3>
<div>
<p>DiskType is the Schema for the disktypes API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>DiskType</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Fleet">Fleet
</h3>
<div>
<p>Fleet is the Schema for the fleets API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Fleet</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetSpec">
FleetSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID identifies the Fleet on provider side.</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.Taint">
[]Taint
</a>
</em>
</td>
<td>
<p>Taints of the Fleet. Only Machines who tolerate all the taints
will land in the Fleet.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetStatus">
FleetStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Instance">Instance
</h3>
<div>
<p>Instance is the Schema for the instances API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Instance</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.InstanceSpec">
InstanceSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>instanceTypeRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>InstanceTypeRef references the instance type of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>fleetSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>FleetSelector selects a suitable Fleet by the given labels.</p>
</td>
</tr>
<tr>
<td>
<code>fleetRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>FleetRef defines the fleet to run the instance in.
If empty, a scheduler will figure out an appropriate pool to run the instance in.</p>
</td>
</tr>
<tr>
<td>
<code>power</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.Power">
Power
</a>
</em>
</td>
<td>
<p>Power is the desired instance power state.
Defaults to PowerOn.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image is the optional URL providing the operating system image of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecret</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>ImagePullSecretRef is an optional secret for pulling the image of a instance.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkInterface">
[]NetworkInterface
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces define a list of network interfaces present on the instance</p>
</td>
</tr>
<tr>
<td>
<code>disks</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.AttachedDisk">
[]AttachedDisk
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Disks are the disks attached to this instance.</p>
</td>
</tr>
<tr>
<td>
<code>ignitionRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.SecretKeySelector">
SecretKeySelector
</a>
</em>
</td>
<td>
<p>IgnitionRef is a reference to a secret containing the ignition YAML for the instance to boot up.
If key is empty, DefaultIgnitionKey will be used as fallback.</p>
</td>
</tr>
<tr>
<td>
<code>efiVars</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.EFIVar">
[]EFIVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EFIVars are variables to pass to EFI while booting up.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.Toleration">
[]Toleration
</a>
</em>
</td>
<td>
<p>Tolerations define tolerations the Instance has. Only fleets whose taints
covered by Tolerations will be considered to run the Instance.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.InstanceStatus">
InstanceStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.InstanceType">InstanceType
</h3>
<div>
<p>InstanceType is the Schema for the instancetypes API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>InstanceType</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>class</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.InstanceTypeClass">
InstanceTypeClass
</a>
</em>
</td>
<td>
<p>Class specifies the class of the InstanceType.
Can either be &lsquo;Continuous&rsquo; or &lsquo;Discrete&rsquo;.</p>
</td>
</tr>
<tr>
<td>
<code>capabilities</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.ResourceList">
ResourceList
</a>
</em>
</td>
<td>
<p>Capabilities are the capabilities of the instance type.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Network">Network
</h3>
<div>
<p>Network is the Schema for the network API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Network</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkStatus">
NetworkStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Subnet">Subnet
</h3>
<div>
<p>Subnet is the Schema for the network API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
core.spheric.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Subnet</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.SubnetSpec">
SubnetSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>networkRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>NetworkRef references the network this subnet is part of.</p>
</td>
</tr>
<tr>
<td>
<code>cidrs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>CIDRs are the primary CIDR ranges of this Subnet.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.SubnetStatus">
SubnetStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.AttachedDisk">AttachedDisk
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>)
</p>
<div>
<p>AttachedDisk defines a disk attached to a instance.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the disk.</p>
</td>
</tr>
<tr>
<td>
<code>device</code><br/>
<em>
string
</em>
</td>
<td>
<p>Device is the device name where the disk should be attached.
Pointer to distinguish between explicit zero and not specified.
If empty, an unused device name will be determined if possible.</p>
</td>
</tr>
<tr>
<td>
<code>AttachedDiskSource</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.AttachedDiskSource">
AttachedDiskSource
</a>
</em>
</td>
<td>
<p>
(Members of <code>AttachedDiskSource</code> are embedded into this type.)
</p>
<p>AttachedDiskSource is the source where the storage for the disk resides at.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.AttachedDiskSource">AttachedDiskSource
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.AttachedDisk">AttachedDisk</a>)
</p>
<div>
<p>AttachedDiskSource specifies the source to use for a disk.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>diskRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>DiskRef instructs to use the specified Disk as source for the attachment.</p>
</td>
</tr>
<tr>
<td>
<code>emptyDisk</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.EmptyDiskSource">
EmptyDiskSource
</a>
</em>
</td>
<td>
<p>EmptyDisk instructs to use a disk offered by the fleet provider.</p>
</td>
</tr>
<tr>
<td>
<code>ephemeral</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.EphemeralDiskSource">
EphemeralDiskSource
</a>
</em>
</td>
<td>
<p>Ephemeral instructs to create an ephemeral (i.e. coupled to the lifetime of the surrounding object)
disk to use.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.AttachedDiskState">AttachedDiskState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.AttachedDiskStatus">AttachedDiskStatus</a>)
</p>
<div>
<p>AttachedDiskState is the infrastructure attachment state a disk can be in.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Attached&#34;</p></td>
<td><p>AttachedDiskStateAttached indicates that a disk has been successfully attached.</p>
</td>
</tr><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>AttachedDiskStatePending indicates that the attachment of a disk is pending.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.AttachedDiskStatus">AttachedDiskStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceStatus">InstanceStatus</a>)
</p>
<div>
<p>AttachedDiskStatus is the status of a disk.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the attached disk.</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.AttachedDiskState">
AttachedDiskState
</a>
</em>
</td>
<td>
<p>State represents the attachment state of a disk.</p>
</td>
</tr>
<tr>
<td>
<code>lastStateTransitionTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>LastStateTransitionTime is the last time the State transitioned.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.ConfigMapKeySelector">ConfigMapKeySelector
</h3>
<div>
<p>ConfigMapKeySelector is a reference to a specific &lsquo;key&rsquo; within a ConfigMap resource.
In some instances, <code>key</code> is a required field.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the referent.
More info: <a href="https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names">https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</a></p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The key of the entry in the ConfigMap resource&rsquo;s <code>data</code> field to be used.
Some instances of this field may be defaulted, in others it may be
required.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DaemonEndpoint">DaemonEndpoint
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetDaemonEndpoints">FleetDaemonEndpoints</a>)
</p>
<div>
<p>DaemonEndpoint contains information about a single Daemon endpoint.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>port</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Port number of the given endpoint.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DiskAccess">DiskAccess
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.DiskStatus">DiskStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>driver</code><br/>
<em>
string
</em>
</td>
<td>
<p>Driver is the name of the drive to use for this volume. Required.</p>
</td>
</tr>
<tr>
<td>
<code>handle</code><br/>
<em>
string
</em>
</td>
<td>
<p>Handle is the unique handle of the volume.</p>
</td>
</tr>
<tr>
<td>
<code>attributes</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Attributes are attributes of the volume to use.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>SecretRef references the (optional) secret containing the data to access the Disk.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DiskSpec">DiskSpec
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Disk">Disk</a>, <a href="#core.spheric.cloud/v1alpha1.DiskTemplateSpec">DiskTemplateSpec</a>)
</p>
<div>
<p>DiskSpec defines the desired state of Disk</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>typeRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>TypeRef references the DiskClass of the Disk.</p>
</td>
</tr>
<tr>
<td>
<code>instanceRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalUIDReference">
LocalUIDReference
</a>
</em>
</td>
<td>
<p>InstanceRef references the using instance of the Disk.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.ResourceList">
ResourceList
</a>
</em>
</td>
<td>
<p>Resources is a description of the Disk&rsquo;s resources and capacity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DiskState">DiskState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.DiskStatus">DiskStatus</a>)
</p>
<div>
<p>DiskState represents the infrastructure state of a Disk.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Available&#34;</p></td>
<td><p>DiskStateAvailable reports whether a Disk is available to be used.</p>
</td>
</tr><tr><td><p>&#34;Error&#34;</p></td>
<td><p>DiskStateError reports that a Disk is in an error state.</p>
</td>
</tr><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>DiskStatePending reports whether a Disk is about to be ready.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DiskStatus">DiskStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Disk">Disk</a>)
</p>
<div>
<p>DiskStatus defines the observed state of Disk</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DiskState">
DiskState
</a>
</em>
</td>
<td>
<p>State represents the infrastructure state of a Disk.</p>
</td>
</tr>
<tr>
<td>
<code>access</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DiskAccess">
DiskAccess
</a>
</em>
</td>
<td>
<p>Access contains information to access the Disk. Must be set when Disk is in DiskStateAvailable.</p>
</td>
</tr>
<tr>
<td>
<code>lastStateTransitionTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>LastStateTransitionTime is the last time the State transitioned between values.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.DiskTemplateSpec">DiskTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.EphemeralDiskSource">EphemeralDiskSource</a>)
</p>
<div>
<p>DiskTemplateSpec is the specification of a Disk template.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DiskSpec">
DiskSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>typeRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>TypeRef references the DiskClass of the Disk.</p>
</td>
</tr>
<tr>
<td>
<code>instanceRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalUIDReference">
LocalUIDReference
</a>
</em>
</td>
<td>
<p>InstanceRef references the using instance of the Disk.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.ResourceList">
ResourceList
</a>
</em>
</td>
<td>
<p>Resources is a description of the Disk&rsquo;s resources and capacity.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.EFIVar">EFIVar
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>)
</p>
<div>
<p>EFIVar is a variable to pass to EFI while booting up.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the EFIVar.</p>
</td>
</tr>
<tr>
<td>
<code>uuid</code><br/>
<em>
string
</em>
</td>
<td>
<p>UUID is the uuid of the EFIVar.</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
<p>Value is the value of the EFIVar.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.EmptyDiskSource">EmptyDiskSource
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.AttachedDiskSource">AttachedDiskSource</a>)
</p>
<div>
<p>EmptyDiskSource is a disk that&rsquo;s offered by the fleet provider.
Usually ephemeral (i.e. deleted when the surrounding entity is deleted), with
varying performance characteristics. Potentially not recoverable.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>sizeLimit</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity">
k8s.io/apimachinery/pkg/api/resource.Quantity
</a>
</em>
</td>
<td>
<p>SizeLimit is the total amount of local storage required for this EmptyDisk disk.
The default is nil which means that the limit is undefined.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.EphemeralDiskSource">EphemeralDiskSource
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.AttachedDiskSource">AttachedDiskSource</a>)
</p>
<div>
<p>EphemeralDiskSource is a definition for an ephemeral (i.e. coupled to the lifetime of the surrounding object)
disk.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>diskTemplate</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DiskTemplateSpec">
DiskTemplateSpec
</a>
</em>
</td>
<td>
<p>DiskTemplate is the template definition of a Disk.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetAddress">FleetAddress
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetStatus">FleetStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetAddressType">
FleetAddressType
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>address</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetAddressType">FleetAddressType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetAddress">FleetAddress</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;ExternalDNS&#34;</p></td>
<td><p>FleetExternalDNS identifies a DNS name which resolves to an IP address which has the characteristics
of FleetExternalIP. The IP it resolves to may or may not be a listed MachineExternalIP address.</p>
</td>
</tr><tr><td><p>&#34;ExternalIP&#34;</p></td>
<td><p>FleetExternalIP identifies an IP address which is, in some way, intended to be more usable from outside
the cluster than an internal IP, though no specific semantics are defined.</p>
</td>
</tr><tr><td><p>&#34;Hostname&#34;</p></td>
<td><p>FleetHostName identifies a name of the fleet. Although every fleet can be assumed
to have a FleetAddress of this type, its exact syntax and semantics are not
defined, and are not consistent between different clusters.</p>
</td>
</tr><tr><td><p>&#34;InternalDNS&#34;</p></td>
<td><p>FleetInternalDNS identifies a DNS name which resolves to an IP address which has
the characteristics of a FleetInternalIP. The IP it resolves to may or may not
be a listed FleetInternalIP address.</p>
</td>
</tr><tr><td><p>&#34;InternalIP&#34;</p></td>
<td><p>FleetInternalIP identifies an IP address which may not be visible to hosts outside the cluster.
By default, it is assumed that apiserver can reach fleet internal IPs, though it is possible
to configure clusters where this is not the case.</p>
<p>FleetInternalIP is the default type of fleet IP, and does not necessarily imply
that the IP is ONLY reachable internally. If a fleet has multiple internal IPs,
no specific semantics are assigned to the additional IPs.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetCondition">FleetCondition
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetStatus">FleetStatus</a>)
</p>
<div>
<p>FleetCondition is one of the conditions of a disk.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetConditionType">
FleetConditionType
</a>
</em>
</td>
<td>
<p>Type is the type of the condition.</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#conditionstatus-v1-core">
Kubernetes core/v1.ConditionStatus
</a>
</em>
</td>
<td>
<p>Status is the status of the condition.</p>
</td>
</tr>
<tr>
<td>
<code>reason</code><br/>
<em>
string
</em>
</td>
<td>
<p>Reason is a machine-readable indication of why the condition is in a certain state.</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message is a human-readable explanation of why the condition has a certain reason / state.</p>
</td>
</tr>
<tr>
<td>
<code>observedGeneration</code><br/>
<em>
int64
</em>
</td>
<td>
<p>ObservedGeneration represents the .metadata.generation that the condition was set based upon.</p>
</td>
</tr>
<tr>
<td>
<code>lastTransitionTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>LastTransitionTime is the last time the status of a condition has transitioned from one state to another.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetConditionType">FleetConditionType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetCondition">FleetCondition</a>)
</p>
<div>
<p>FleetConditionType is a type a FleetCondition can have.</p>
</div>
<h3 id="core.spheric.cloud/v1alpha1.FleetDaemonEndpoints">FleetDaemonEndpoints
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetStatus">FleetStatus</a>)
</p>
<div>
<p>FleetDaemonEndpoints lists ports opened by daemons running on the Fleet.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>sphereletEndpoint</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.DaemonEndpoint">
DaemonEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Endpoint on which spherelet is listening.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetSpec">FleetSpec
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Fleet">Fleet</a>)
</p>
<div>
<p>FleetSpec defines the desired state of Fleet</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID identifies the Fleet on provider side.</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.Taint">
[]Taint
</a>
</em>
</td>
<td>
<p>Taints of the Fleet. Only Machines who tolerate all the taints
will land in the Fleet.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetState">FleetState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetStatus">FleetStatus</a>)
</p>
<div>
<p>FleetState is a state a Fleet can be in.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Error&#34;</p></td>
<td><p>FleetStateError marks a Fleet in an error state.</p>
</td>
</tr><tr><td><p>&#34;Offline&#34;</p></td>
<td><p>FleetStateOffline marks a Fleet as offline.</p>
</td>
</tr><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>FleetStatePending marks a Fleet as pending readiness.</p>
</td>
</tr><tr><td><p>&#34;Ready&#34;</p></td>
<td><p>FleetStateReady marks a Fleet as ready for accepting a Machine.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.FleetStatus">FleetStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Fleet">Fleet</a>)
</p>
<div>
<p>FleetStatus defines the observed state of Fleet</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetState">
FleetState
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetCondition">
[]FleetCondition
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetAddress">
[]FleetAddress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>daemonEndpoints</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.FleetDaemonEndpoints">
FleetDaemonEndpoints
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>capacity</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.ResourceList">
ResourceList
</a>
</em>
</td>
<td>
<p>Capacity represents the total resources of a fleet.</p>
</td>
</tr>
<tr>
<td>
<code>allocatable</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.ResourceList">
ResourceList
</a>
</em>
</td>
<td>
<p>Allocatable represents the resources of a fleet that are available for scheduling.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.InstanceExecOptions">InstanceExecOptions
</h3>
<div>
<p>InstanceExecOptions is the query options to a Instance&rsquo;s remote exec call</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>insecureSkipTLSVerifyBackend</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Instance">Instance</a>)
</p>
<div>
<p>InstanceSpec defines the desired state of Instance</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceTypeRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>InstanceTypeRef references the instance type of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>fleetSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>FleetSelector selects a suitable Fleet by the given labels.</p>
</td>
</tr>
<tr>
<td>
<code>fleetRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>FleetRef defines the fleet to run the instance in.
If empty, a scheduler will figure out an appropriate pool to run the instance in.</p>
</td>
</tr>
<tr>
<td>
<code>power</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.Power">
Power
</a>
</em>
</td>
<td>
<p>Power is the desired instance power state.
Defaults to PowerOn.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image is the optional URL providing the operating system image of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecret</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>ImagePullSecretRef is an optional secret for pulling the image of a instance.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkInterface">
[]NetworkInterface
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces define a list of network interfaces present on the instance</p>
</td>
</tr>
<tr>
<td>
<code>disks</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.AttachedDisk">
[]AttachedDisk
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Disks are the disks attached to this instance.</p>
</td>
</tr>
<tr>
<td>
<code>ignitionRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.SecretKeySelector">
SecretKeySelector
</a>
</em>
</td>
<td>
<p>IgnitionRef is a reference to a secret containing the ignition YAML for the instance to boot up.
If key is empty, DefaultIgnitionKey will be used as fallback.</p>
</td>
</tr>
<tr>
<td>
<code>efiVars</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.EFIVar">
[]EFIVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EFIVars are variables to pass to EFI while booting up.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.Toleration">
[]Toleration
</a>
</em>
</td>
<td>
<p>Tolerations define tolerations the Instance has. Only fleets whose taints
covered by Tolerations will be considered to run the Instance.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.InstanceState">InstanceState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceStatus">InstanceStatus</a>)
</p>
<div>
<p>InstanceState is the state of a instance.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>InstanceStatePending means the Instance has been accepted by the system, but not yet completely started.
This includes time before being bound to a Fleet, as well as time spent setting up the Instance on that
Fleet.</p>
</td>
</tr><tr><td><p>&#34;Running&#34;</p></td>
<td><p>InstanceStateRunning means the instance is running on a Fleet.</p>
</td>
</tr><tr><td><p>&#34;Shutdown&#34;</p></td>
<td><p>InstanceStateShutdown means the instance is shut down.</p>
</td>
</tr><tr><td><p>&#34;Terminated&#34;</p></td>
<td><p>InstanceStateTerminated means the instance has been permanently stopped and cannot be started.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.InstanceStatus">InstanceStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Instance">Instance</a>)
</p>
<div>
<p>InstanceStatus defines the observed state of Instance</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the provider specific instance ID in the format &lsquo;<type>://<instance_id>&rsquo;.</p>
</td>
</tr>
<tr>
<td>
<code>observedGeneration</code><br/>
<em>
int64
</em>
</td>
<td>
<p>ObservedGeneration is the last generation the Fleet observed of the Instance.</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.InstanceState">
InstanceState
</a>
</em>
</td>
<td>
<p>State is the infrastructure state of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkInterfaceStatus">
[]NetworkInterfaceStatus
</a>
</em>
</td>
<td>
<p>NetworkInterfaces is the list of network interface states for the instance.</p>
</td>
</tr>
<tr>
<td>
<code>disks</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.AttachedDiskStatus">
[]AttachedDiskStatus
</a>
</em>
</td>
<td>
<p>Disks is the list of disk states for the instance.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.InstanceTypeClass">InstanceTypeClass
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceType">InstanceType</a>)
</p>
<div>
<p>InstanceTypeClass denotes the type of InstanceType.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Continuous&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;Discrete&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.LocalObjectReference">LocalObjectReference
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.AttachedDiskSource">AttachedDiskSource</a>, <a href="#core.spheric.cloud/v1alpha1.DiskAccess">DiskAccess</a>, <a href="#core.spheric.cloud/v1alpha1.DiskSpec">DiskSpec</a>, <a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>, <a href="#core.spheric.cloud/v1alpha1.SubnetSpec">SubnetSpec</a>)
</p>
<div>
<p>LocalObjectReference contains enough information to let you locate the
referenced object inside the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the referent.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.LocalUIDReference">LocalUIDReference
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.DiskSpec">DiskSpec</a>)
</p>
<div>
<p>LocalUIDReference is a reference to another entity including its UID</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the referenced entity.</p>
</td>
</tr>
<tr>
<td>
<code>uid</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/types#UID">
k8s.io/apimachinery/pkg/types.UID
</a>
</em>
</td>
<td>
<p>UID is the UID of the referenced entity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.NetworkInterface">NetworkInterface
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>)
</p>
<div>
<p>NetworkInterface is the definition of a single interface</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the network interface.</p>
</td>
</tr>
<tr>
<td>
<code>subnetRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.SubnetReference">
SubnetReference
</a>
</em>
</td>
<td>
<p>SubnetRef references the Subnet this NetworkInterface is connected to</p>
</td>
</tr>
<tr>
<td>
<code>ipFamilies</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#ipfamily-v1-core">
[]Kubernetes core/v1.IPFamily
</a>
</em>
</td>
<td>
<p>IPFamilies defines which IPFamilies this NetworkInterface is supporting</p>
</td>
</tr>
<tr>
<td>
<code>ips</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>IPs are the literal requested IPs for this NetworkInterface.</p>
</td>
</tr>
<tr>
<td>
<code>accessIPFamilies</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#ipfamily-v1-core">
[]Kubernetes core/v1.IPFamily
</a>
</em>
</td>
<td>
<p>AccessIPFamilies are the access configuration IP families.</p>
</td>
</tr>
<tr>
<td>
<code>accessIPs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AccessIPs are the literal request access IPs.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.NetworkInterfaceState">NetworkInterfaceState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.NetworkInterfaceStatus">NetworkInterfaceStatus</a>)
</p>
<div>
<p>NetworkInterfaceState is the infrastructure attachment state a NetworkInterface can be in.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Attached&#34;</p></td>
<td><p>NetworkInterfaceStateAttached indicates that a network interface has been successfully attached.</p>
</td>
</tr><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>NetworkInterfaceStatePending indicates that the attachment of a network interface is pending.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.NetworkInterfaceStatus">NetworkInterfaceStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceStatus">InstanceStatus</a>)
</p>
<div>
<p>NetworkInterfaceStatus reports the status of an NetworkInterfaceSource.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the NetworkInterface to whom the status belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>ips</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>IPs are the ips allocated for the network interface.</p>
</td>
</tr>
<tr>
<td>
<code>accessIPs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AccessIPs are the allocated access IPs for the network interface.</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkInterfaceState">
NetworkInterfaceState
</a>
</em>
</td>
<td>
<p>State represents the attachment state of a NetworkInterface.</p>
</td>
</tr>
<tr>
<td>
<code>lastStateTransitionTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>LastStateTransitionTime is the last time the State transitioned.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.NetworkSpec">NetworkSpec
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Network">Network</a>)
</p>
<div>
<p>NetworkSpec defines the desired state of Network</p>
</div>
<h3 id="core.spheric.cloud/v1alpha1.NetworkState">NetworkState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.NetworkStatus">NetworkStatus</a>)
</p>
<div>
<p>NetworkState is the state of a network.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Available&#34;</p></td>
<td><p>NetworkStateAvailable means the network is ready to use.</p>
</td>
</tr><tr><td><p>&#34;Error&#34;</p></td>
<td><p>NetworkStateError means the network is in an error state.</p>
</td>
</tr><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>NetworkStatePending means the network is being provisioned.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.NetworkStatus">NetworkStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Network">Network</a>)
</p>
<div>
<p>NetworkStatus defines the observed state of Network</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.NetworkState">
NetworkState
</a>
</em>
</td>
<td>
<p>State is the state of the machine.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.ObjectSelector">ObjectSelector
</h3>
<div>
<p>ObjectSelector specifies how to select objects of a certain kind.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>kind</code><br/>
<em>
string
</em>
</td>
<td>
<p>Kind is the kind of object to select.</p>
</td>
</tr>
<tr>
<td>
<code>LabelSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>
(Members of <code>LabelSelector</code> are embedded into this type.)
</p>
<p>LabelSelector is the label selector to select objects of the specified Kind by.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Power">Power
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>)
</p>
<div>
<p>Power is the desired power state of a Instance.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Off&#34;</p></td>
<td><p>PowerOff indicates that a Instance should be powered off.</p>
</td>
</tr><tr><td><p>&#34;On&#34;</p></td>
<td><p>PowerOn indicates that a Instance should be powered on.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.ResourceName">ResourceName
(<code>string</code> alias)</h3>
<div>
<p>ResourceName is the name of a resource, most often used alongside a resource.Quantity.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;cpu&#34;</p></td>
<td><p>ResourceCPU is the amount of cpu in cores.</p>
</td>
</tr><tr><td><p>&#34;iops&#34;</p></td>
<td><p>ResourceIOPS defines max IOPS in input/output operations per second.</p>
</td>
</tr><tr><td><p>&#34;memory&#34;</p></td>
<td><p>ResourceMemory is the amount of memory in bytes.</p>
</td>
</tr><tr><td><p>&#34;storage&#34;</p></td>
<td><p>ResourceStorage is the amount of storage, in bytes.</p>
</td>
</tr><tr><td><p>&#34;tps&#34;</p></td>
<td><p>ResourceTPS defines max throughput per second. (e.g. 1Gi)</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.SecretKeySelector">SecretKeySelector
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>)
</p>
<div>
<p>SecretKeySelector is a reference to a specific &lsquo;key&rsquo; within a Secret resource.
In some instances, <code>key</code> is a required field.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the referent.
More info: <a href="https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names">https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</a></p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The key of the entry in the Secret resource&rsquo;s <code>data</code> field to be used.
Some instances of this field may be defaulted, in others it may be
required.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.SubnetReference">SubnetReference
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.NetworkInterface">NetworkInterface</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>networkName</code><br/>
<em>
string
</em>
</td>
<td>
<p>NetworkName is the name of the referenced network.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the referenced subnet.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.SubnetSpec">SubnetSpec
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Subnet">Subnet</a>)
</p>
<div>
<p>SubnetSpec defines the desired state of Subnet</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>networkRef</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<p>NetworkRef references the network this subnet is part of.</p>
</td>
</tr>
<tr>
<td>
<code>cidrs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>CIDRs are the primary CIDR ranges of this Subnet.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.SubnetState">SubnetState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.SubnetStatus">SubnetStatus</a>)
</p>
<div>
<p>SubnetState is the state of a network.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Available&#34;</p></td>
<td><p>SubnetStateAvailable means the network is ready to use.</p>
</td>
</tr><tr><td><p>&#34;Error&#34;</p></td>
<td><p>SubnetStateError means the network is in an error state.</p>
</td>
</tr><tr><td><p>&#34;Pending&#34;</p></td>
<td><p>SubnetStatePending means the network is being provisioned.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.SubnetStatus">SubnetStatus
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Subnet">Subnet</a>)
</p>
<div>
<p>SubnetStatus defines the observed state of Subnet</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>state</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.SubnetState">
SubnetState
</a>
</em>
</td>
<td>
<p>State is the state of the machine.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Taint">Taint
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.FleetSpec">FleetSpec</a>)
</p>
<div>
<p>Taint marks an effect with a value on a target resource pool.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>The taint key to be applied to a resource pool.</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
<p>The taint value corresponding to the taint key.</p>
</td>
</tr>
<tr>
<td>
<code>effect</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.TaintEffect">
TaintEffect
</a>
</em>
</td>
<td>
<p>The effect of the taint on resources
that do not tolerate the taint.
Valid effects are NoSchedule, PreferNoSchedule and NoExecute.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.TaintEffect">TaintEffect
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Taint">Taint</a>, <a href="#core.spheric.cloud/v1alpha1.Toleration">Toleration</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;NoSchedule&#34;</p></td>
<td><p>TaintEffectNoSchedule causes not to allow new resources to schedule onto the resource pool unless they tolerate
the taint, but allow all already-running resources to continue running.
Enforced by the scheduler.</p>
</td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.Toleration">Toleration
</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.InstanceSpec">InstanceSpec</a>)
</p>
<div>
<p>Toleration marks the resource the toleration is attached to tolerate any taint that matches
the triple <key,value,effect> using the matching operator <operator>.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>Key is the taint key that the toleration applies to. Empty means match all taint keys.
If the key is empty, operator must be Exists; this combination means to match all values and all keys.</p>
</td>
</tr>
<tr>
<td>
<code>operator</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.TolerationOperator">
TolerationOperator
</a>
</em>
</td>
<td>
<p>Operator represents a key&rsquo;s relationship to the value.
Valid operators are Exists and Equal. Defaults to Equal.
Exists is equivalent to wildcard for value, so that a resource can
tolerate all taints of a particular category.</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
<p>Value is the taint value the toleration matches to.
If the operator is Exists, the value should be empty, otherwise just a regular string.</p>
</td>
</tr>
<tr>
<td>
<code>effect</code><br/>
<em>
<a href="#core.spheric.cloud/v1alpha1.TaintEffect">
TaintEffect
</a>
</em>
</td>
<td>
<p>Effect indicates the taint effect to match. Empty means match all taint effects.
When specified, allowed values are NoSchedule.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.TolerationOperator">TolerationOperator
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#core.spheric.cloud/v1alpha1.Toleration">Toleration</a>)
</p>
<div>
<p>TolerationOperator is the set of operators that can be used in a toleration.</p>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Equal&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;Exists&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="core.spheric.cloud/v1alpha1.UIDReference">UIDReference
</h3>
<div>
<p>UIDReference is a reference to another entity in a potentially different namespace including its UID.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>namespace</code><br/>
<em>
string
</em>
</td>
<td>
<p>Namespace is the namespace of the referenced entity. If empty,
the same namespace as the referring resource is implied.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the referenced entity.</p>
</td>
</tr>
<tr>
<td>
<code>uid</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/types#UID">
k8s.io/apimachinery/pkg/types.UID
</a>
</em>
</td>
<td>
<p>UID is the UID of the referenced entity.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
</em></p>
