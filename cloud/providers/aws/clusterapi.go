/*
Copyright The Pharmer Authors.

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
package aws

func (cm *ClusterManager) GetClusterAPIComponents() (string, error) {
	return ControllerManager, nil
}

// ControllerManager defines the clusterapi crds and controllers
const ControllerManager = `
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: awsclusterproviderspecs.awsprovider.k8s.io
spec:
  group: awsprovider.k8s.io
  names:
    kind: AWSClusterProviderSpec
    plural: awsclusterproviderspecs
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        caKeyPair:
          description: CAKeyPair is the key pair for ca certs.
          properties:
            cert:
              description: base64 encoded cert and key
              format: byte
              type: string
            key:
              format: byte
              type: string
          required:
          - cert
          - key
          type: object
        etcdCAKeyPair:
          description: EtcdCAKeyPair is the key pair for etcd.
          properties:
            cert:
              description: base64 encoded cert and key
              format: byte
              type: string
            key:
              format: byte
              type: string
          required:
          - cert
          - key
          type: object
        frontProxyCAKeyPair:
          description: FrontProxyCAKeyPair is the key pair for FrontProxyKeyPair.
          properties:
            cert:
              description: base64 encoded cert and key
              format: byte
              type: string
            key:
              format: byte
              type: string
          required:
          - cert
          - key
          type: object
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        networkSpec:
          description: NetworkSpec encapsulates all things related to AWS network.
          properties:
            subnets:
              description: Subnets configuration.
              items:
                properties:
                  availabilityZone:
                    description: AvailabilityZone defines the availability zone to
                      use for this subnet in the cluster's region.
                    type: string
                  cidrBlock:
                    description: CidrBlock is the CIDR block to be used when the provider
                      creates a managed VPC.
                    type: string
                  id:
                    description: ID defines a unique identifier to reference this
                      resource.
                    type: string
                  isPublic:
                    description: IsPublic defines the subnet as a public subnet. Refer
                      to the AWS documentation for further information.
                    type: boolean
                  natGatewayId:
                    description: NatGatewayID is the NAT gateway id associated with
                      the subnet. Ignored if the subnet is public.
                    type: string
                  routeTableId:
                    description: RouteTableID is the routing table id associated with
                      the subnet.
                    type: string
                  tags:
                    description: Tags is a collection of tags describing the resource.
                    type: object
                required:
                - routeTableId
                - natGatewayId
                type: object
              type: array
            vpc:
              description: VPC configuration.
              properties:
                cidrBlock:
                  description: CidrBlock is the CIDR block to be used when the provider
                    creates a managed VPC. Defaults to 10.0.0.0/16.
                  type: string
                id:
                  description: ID is the vpc-id of the VPC this provider should use
                    to create resources.
                  type: string
                internetGatewayId:
                  description: InternetGatewayID is the id of the internet gateway
                    associated with the VPC.
                  type: string
                tags:
                  description: Tags is a collection of tags describing the resource.
                  type: object
              type: object
          type: object
        region:
          description: The AWS Region the cluster lives in.
          type: string
        saKeyPair:
          description: SAKeyPair is the service account key pair.
          properties:
            cert:
              description: base64 encoded cert and key
              format: byte
              type: string
            key:
              format: byte
              type: string
          required:
          - cert
          - key
          type: object
        sshKeyName:
          description: SSHKeyName is the name of the ssh key to attach to the bastion
            host.
          type: string
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: awsclusterproviderstatuses.awsprovider.k8s.io
spec:
  group: awsprovider.k8s.io
  names:
    kind: AWSClusterProviderStatus
    plural: awsclusterproviderstatuses
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        bastion:
          properties:
            ebsOptimized:
              description: Indicates whether the instance is optimized for Amazon
                EBS I/O.
              type: boolean
            enaSupport:
              description: Specifies whether enhanced networking with ENA is enabled.
              type: boolean
            iamProfile:
              description: The name of the IAM instance profile associated with the
                instance, if applicable.
              type: string
            id:
              type: string
            imageId:
              description: The ID of the AMI used to launch the instance.
              type: string
            instanceState:
              description: The current state of the instance.
              type: string
            keyName:
              description: The name of the SSH key pair.
              type: string
            privateIp:
              description: The private IPv4 address assigned to the instance.
              type: string
            publicIp:
              description: The public IPv4 address assigned to the instance, if applicable.
              type: string
            securityGroupIds:
              description: SecurityGroupIDs are one or more security group IDs this
                instance belongs to.
              items:
                type: string
              type: array
            subnetId:
              description: The ID of the subnet of the instance.
              type: string
            tags:
              description: The tags associated with the instance.
              type: object
            type:
              description: The instance type.
              type: string
            userData:
              description: UserData is the raw data script passed to the instance
                which is run upon bootstrap. This field must not be base64 encoded
                and should only be used when running a new instance.
              type: string
          required:
          - id
          type: object
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        network:
          properties:
            apiServerElb:
              description: APIServerELB is the Kubernetes api server classic load
                balancer.
              properties:
                attributes:
                  description: Attributes defines extra attributes associated with
                    the load balancer.
                  properties:
                    idleTimeout:
                      description: IdleTimeout is time that the connection is allowed
                        to be idle (no data has been sent over the connection) before
                        it is closed by the load balancer.
                      format: int64
                      type: integer
                  type: object
                dnsName:
                  description: DNSName is the dns name of the load balancer.
                  type: string
                healthChecks:
                  description: HealthCheck is the classic elb health check associated
                    with the load balancer.
                  properties:
                    healthyThreshold:
                      format: int64
                      type: integer
                    interval:
                      format: int64
                      type: integer
                    target:
                      type: string
                    timeout:
                      format: int64
                      type: integer
                    unhealthyThreshold:
                      format: int64
                      type: integer
                  required:
                  - target
                  - interval
                  - timeout
                  - healthyThreshold
                  - unhealthyThreshold
                  type: object
                listeners:
                  description: Listeners is an array of classic elb listeners associated
                    with the load balancer. There must be at least one.
                  items:
                    properties:
                      instancePort:
                        format: int64
                        type: integer
                      instanceProtocol:
                        type: string
                      port:
                        format: int64
                        type: integer
                      protocol:
                        type: string
                    required:
                    - protocol
                    - port
                    - instanceProtocol
                    - instancePort
                    type: object
                  type: array
                name:
                  description: The name of the load balancer. It must be unique within
                    the set of load balancers defined in the region. It also serves
                    as identifier.
                  type: string
                scheme:
                  description: Scheme is the load balancer scheme, either internet-facing
                    or private.
                  type: string
                securityGroupIds:
                  description: SecurityGroupIDs is an array of security groups assigned
                    to the load balancer.
                  items:
                    type: string
                  type: array
                subnetIds:
                  description: SubnetIDs is an array of subnets in the VPC attached
                    to the load balancer.
                  items:
                    type: string
                  type: array
                tags:
                  description: Tags is a map of tags associated with the load balancer.
                  type: object
              type: object
            securityGroups:
              description: SecurityGroups is a map from the role/kind of the security
                group to its unique name, if any.
              type: object
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: awsmachineproviderspecs.awsprovider.k8s.io
spec:
  group: awsprovider.k8s.io
  names:
    kind: AWSMachineProviderSpec
    plural: awsmachineproviderspecs
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        additionalSecurityGroups:
          description: AdditionalSecurityGroups is an array of references to security
            groups that should be applied to the instance. These security groups would
            be set in addition to any security groups defined at the cluster level
            or in the actuator.
          items:
            properties:
              arn:
                description: ARN of resource
                type: string
              filters:
                description: 'Filters is a set of key/value pairs used to identify
                  a resource They are applied according to the rules defined by the
                  AWS API: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html'
                items:
                  properties:
                    name:
                      description: Name of the filter. Filter names are case-sensitive.
                      type: string
                    values:
                      description: Values includes one or more filter values. Filter
                        values are case-sensitive.
                      items:
                        type: string
                      type: array
                  required:
                  - name
                  - values
                  type: object
                type: array
              id:
                description: ID of resource
                type: string
            type: object
          type: array
        additionalTags:
          description: AdditionalTags is the set of tags to add to an instance, in
            addition to the ones added by default by the actuator. These tags are
            additive. The actuator will ensure these tags are present, but will not
            remove any other tags that may exist on the instance.
          type: object
        ami:
          description: AMI is the reference to the AMI from which to create the machine
            instance.
          properties:
            arn:
              description: ARN of resource
              type: string
            filters:
              description: 'Filters is a set of key/value pairs used to identify a
                resource They are applied according to the rules defined by the AWS
                API: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html'
              items:
                properties:
                  name:
                    description: Name of the filter. Filter names are case-sensitive.
                    type: string
                  values:
                    description: Values includes one or more filter values. Filter
                      values are case-sensitive.
                    items:
                      type: string
                    type: array
                required:
                - name
                - values
                type: object
              type: array
            id:
              description: ID of resource
              type: string
          type: object
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        iamInstanceProfile:
          description: IAMInstanceProfile is a name of an IAM instance profile to
            assign to the instance
          type: string
        instanceType:
          description: 'InstanceType is the type of instance to create. Example: m4.xlarge'
          type: string
        keyName:
          description: KeyName is the name of the SSH key to install on the instance.
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        publicIP:
          description: 'PublicIP specifies whether the instance should get a public
            IP. Precedence for this setting is as follows: 1. This field if set 2.
            Cluster/flavor setting 3. Subnet default'
          type: boolean
        subnet:
          description: Subnet is a reference to the subnet to use for this instance.
            If not specified, the cluster subnet will be used.
          properties:
            arn:
              description: ARN of resource
              type: string
            filters:
              description: 'Filters is a set of key/value pairs used to identify a
                resource They are applied according to the rules defined by the AWS
                API: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html'
              items:
                properties:
                  name:
                    description: Name of the filter. Filter names are case-sensitive.
                    type: string
                  values:
                    description: Values includes one or more filter values. Filter
                      values are case-sensitive.
                    items:
                      type: string
                    type: array
                required:
                - name
                - values
                type: object
              type: array
            id:
              description: ID of resource
              type: string
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: awsmachineproviderstatuses.awsprovider.k8s.io
spec:
  group: awsprovider.k8s.io
  names:
    kind: AWSMachineProviderStatus
    plural: awsmachineproviderstatuses
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        conditions:
          description: Conditions is a set of conditions associated with the Machine
            to indicate errors or other status
          items:
            properties:
              lastProbeTime:
                description: LastProbeTime is the last time we probed the condition.
                format: date-time
                type: string
              lastTransitionTime:
                description: LastTransitionTime is the last time the condition transitioned
                  from one status to another.
                format: date-time
                type: string
              message:
                description: Message is a human-readable message indicating details
                  about last transition.
                type: string
              reason:
                description: Reason is a unique, one-word, CamelCase reason for the
                  condition's last transition.
                type: string
              status:
                description: Status is the status of the condition.
                type: string
              type:
                description: Type is the type of the condition.
                type: string
            required:
            - type
            - status
            - lastProbeTime
            - lastTransitionTime
            - reason
            - message
            type: object
          type: array
        instanceID:
          description: InstanceID is the instance ID of the machine created in AWS
          type: string
        instanceState:
          description: InstanceState is the state of the AWS instance for this machine
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: aws-provider-manager-role
rules:
- apiGroups:
  - awsprovider.k8s.io
  resources:
  - awsclusterproviderconfigs
  - awsclusterproviderstatuses
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - clusters
  - clusters/status
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - awsprovider.k8s.io
  resources:
  - awsmachineproviderconfigs
  - awsmachineproviderstatuses
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - machines
  - machines/status
  - machinedeployments
  - machinedeployments/status
  - machinesets
  - machinesets/status
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - clusters
  - clusters/status
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes
  - events
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: null
  name: aws-provider-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: aws-provider-manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: aws-provider-system
---
apiVersion: v1
kind: Service
metadata:
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: aws-provider-controller-manager-service
  namespace: aws-provider-system
spec:
  ports:
  - port: 443
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: aws-provider-controller-manager
  namespace: aws-provider-system
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
      controller-tools.k8s.io: "1.0"
  serviceName: aws-provider-controller-manager-service
  template:
    metadata:
      labels:
        control-plane: controller-manager
        controller-tools.k8s.io: "1.0"
    spec:
      containers:
      - args:
        - -v=3
        - -logtostderr=true
        - -stderrthreshold=INFO
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        image: pharmer/cluster-api-provider-aws:0.0.1
        imagePullPolicy: IfNotPresent
        name: manager
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
        volumeMounts:
        - mountPath: /etc/kubernetes
          name: config
        - mountPath: /etc/ssl/certs
          name: certs
        - mountPath: /root/.aws
          name: credentials
      terminationGracePeriodSeconds: 10
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoExecute
        key: node.alpha.kubernetes.io/notReady
        operator: Exists
      - effect: NoExecute
        key: node.alpha.kubernetes.io/unreachable
        operator: Exists
      volumes:
      - hostPath:
          path: /etc/kubernetes
        name: config
      - hostPath:
          path: /etc/ssl/certs
        name: certs
      - name: credentials
        secret:
          secretName: aws-provider-manager-bootstrap-credentials
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: cluster-api-system
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: clusters.cluster.k8s.io
spec:
  group: cluster.k8s.io
  names:
    kind: Cluster
    plural: clusters
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            clusterNetwork:
              description: Cluster network configuration
              properties:
                pods:
                  description: The network ranges from which Pod networks are allocated.
                  properties:
                    cidrBlocks:
                      items:
                        type: string
                      type: array
                  required:
                  - cidrBlocks
                  type: object
                serviceDomain:
                  description: Domain name for services.
                  type: string
                services:
                  description: The network ranges from which service VIPs are allocated.
                  properties:
                    cidrBlocks:
                      items:
                        type: string
                      type: array
                  required:
                  - cidrBlocks
                  type: object
              required:
              - services
              - pods
              - serviceDomain
              type: object
            providerSpec:
              description: Provider-specific serialized configuration to use during
                cluster creation. It is recommended that providers maintain their
                own versioned API types that should be serialized/deserialized from
                this field.
              properties:
                value:
                  description: Value is an inlined, serialized representation of the
                    resource configuration. It is recommended that providers maintain
                    their own versioned API types that should be serialized/deserialized
                    from this field, akin to component config.
                  type: object
                valueFrom:
                  description: Source for the provider configuration. Cannot be used
                    if value is not empty.
                  properties:
                    machineClass:
                      description: The machine class from which the provider config
                        should be sourced.
                      properties:
                        provider:
                          description: Provider is the name of the cloud-provider
                            which MachineClass is intended for.
                          type: string
                      type: object
                  type: object
              type: object
          required:
          - clusterNetwork
          type: object
        status:
          properties:
            apiEndpoints:
              description: APIEndpoint represents the endpoint to communicate with
                the IP.
              items:
                properties:
                  host:
                    description: The hostname on which the API server is serving.
                    type: string
                  port:
                    description: The port on which the API server is serving.
                    format: int64
                    type: integer
                required:
                - host
                - port
                type: object
              type: array
            errorMessage:
              description: If set, indicates that there is a problem reconciling the
                state, and will be set to a descriptive error message.
              type: string
            errorReason:
              description: If set, indicates that there is a problem reconciling the
                state, and will be set to a token value suitable for programmatic
                interpretation.
              type: string
            providerStatus:
              description: Provider-specific status. It is recommended that providers
                maintain their own versioned API types that should be serialized/deserialized
                from this field.
              type: object
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: machineclasses.cluster.k8s.io
spec:
  group: cluster.k8s.io
  names:
    kind: MachineClass
    plural: machineclasses
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        providerSpec:
          description: Provider-specific configuration to use during node creation.
          type: object
      required:
      - providerSpec
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: machinedeployments.cluster.k8s.io
spec:
  group: cluster.k8s.io
  names:
    kind: MachineDeployment
    plural: machinedeployments
  scope: Namespaced
  subresources:
    scale:
      labelSelectorPath: .status.labelSelector
      specReplicasPath: .spec.replicas
      statusReplicasPath: .status.replicas
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            minReadySeconds:
              description: Minimum number of seconds for which a newly created machine
                should be ready. Defaults to 0 (machine will be considered available
                as soon as it is ready)
              format: int32
              type: integer
            paused:
              description: Indicates that the deployment is paused.
              type: boolean
            progressDeadlineSeconds:
              description: The maximum time in seconds for a deployment to make progress
                before it is considered to be failed. The deployment controller will
                continue to process failed deployments and a condition with a ProgressDeadlineExceeded
                reason will be surfaced in the deployment status. Note that progress
                will not be estimated during the time a deployment is paused. Defaults
                to 600s.
              format: int32
              type: integer
            replicas:
              description: Number of desired machines. Defaults to 1. This is a pointer
                to distinguish between explicit zero and not specified.
              format: int32
              type: integer
            revisionHistoryLimit:
              description: The number of old MachineSets to retain to allow rollback.
                This is a pointer to distinguish between explicit zero and not specified.
                Defaults to 1.
              format: int32
              type: integer
            selector:
              description: Label selector for machines. Existing MachineSets whose
                machines are selected by this will be the ones affected by this deployment.
                It must match the machine template's labels.
              type: object
            strategy:
              description: The deployment strategy to use to replace existing machines
                with new ones.
              properties:
                rollingUpdate:
                  description: Rolling update config params. Present only if MachineDeploymentStrategyType
                    = RollingUpdate.
                  properties:
                    maxSurge:
                      description: 'The maximum number of machines that can be scheduled
                        above the desired number of machines. Value can be an absolute
                        number (ex: 5) or a percentage of desired machines (ex: 10%).
                        This can not be 0 if MaxUnavailable is 0. Absolute number
                        is calculated from percentage by rounding up. Defaults to
                        1. Example: when this is set to 30%, the new MachineSet can
                        be scaled up immediately when the rolling update starts, such
                        that the total number of old and new machines do not exceed
                        130% of desired machines. Once old machines have been killed,
                        new MachineSet can be scaled up further, ensuring that total
                        number of machines running at any time during the update is
                        at most 130% of desired machines.'
                      oneOf:
                      - type: string
                      - type: integer
                    maxUnavailable:
                      description: 'The maximum number of machines that can be unavailable
                        during the update. Value can be an absolute number (ex: 5)
                        or a percentage of desired machines (ex: 10%). Absolute number
                        is calculated from percentage by rounding down. This can not
                        be 0 if MaxSurge is 0. Defaults to 0. Example: when this is
                        set to 30%, the old MachineSet can be scaled down to 70% of
                        desired machines immediately when the rolling update starts.
                        Once new machines are ready, old MachineSet can be scaled
                        down further, followed by scaling up the new MachineSet, ensuring
                        that the total number of machines available at all times during
                        the update is at least 70% of desired machines.'
                      oneOf:
                      - type: string
                      - type: integer
                  type: object
                type:
                  description: Type of deployment. Currently the only supported strategy
                    is "RollingUpdate". Default is RollingUpdate.
                  type: string
              type: object
            template:
              description: Template describes the machines that will be created.
              properties:
                metadata:
                  description: 'Standard object''s metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata'
                  type: object
                spec:
                  description: 'Specification of the desired behavior of the machine.
                    More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status'
                  properties:
                    configSource:
                      description: ConfigSource is used to populate in the associated
                        Node for dynamic kubelet config. This field already exists
                        in Node, so any updates to it in the Machine spec will be
                        automatically copied to the linked NodeRef from the status.
                        The rest of dynamic kubelet config support should then work
                        as-is.
                      type: object
                    metadata:
                      description: ObjectMeta will autopopulate the Node created.
                        Use this to indicate what labels, annotations, name prefix,
                        etc., should be used when creating the Node.
                      type: object
                    providerSpec:
                      description: ProviderSpec details Provider-specific configuration
                        to use during node creation.
                      properties:
                        value:
                          description: Value is an inlined, serialized representation
                            of the resource configuration. It is recommended that
                            providers maintain their own versioned API types that
                            should be serialized/deserialized from this field, akin
                            to component config.
                          type: object
                        valueFrom:
                          description: Source for the provider configuration. Cannot
                            be used if value is not empty.
                          properties:
                            machineClass:
                              description: The machine class from which the provider
                                config should be sourced.
                              properties:
                                provider:
                                  description: Provider is the name of the cloud-provider
                                    which MachineClass is intended for.
                                  type: string
                              type: object
                          type: object
                      type: object
                    taints:
                      description: Taints is the full, authoritative list of taints
                        to apply to the corresponding Node. This list will overwrite
                        any modifications made to the Node on an ongoing basis.
                      items:
                        type: object
                      type: array
                    versions:
                      description: Versions of key software to use. This field is
                        optional at cluster creation time, and omitting the field
                        indicates that the cluster installation tool should select
                        defaults for the user. These defaults may differ based on
                        the cluster installer, but the tool should populate the values
                        it uses when persisting Machine objects. A Machine spec missing
                        this field at runtime is invalid.
                      properties:
                        controlPlane:
                          description: ControlPlane is the semantic version of the
                            Kubernetes control plane to run. This should only be populated
                            when the machine is a control plane.
                          type: string
                        kubelet:
                          description: Kubelet is the semantic version of kubelet
                            to run
                          type: string
                      required:
                      - kubelet
                      type: object
                  required:
                  - providerSpec
                  type: object
              type: object
          required:
          - selector
          - template
          type: object
        status:
          properties:
            availableReplicas:
              description: Total number of available machines (ready for at least
                minReadySeconds) targeted by this deployment.
              format: int32
              type: integer
            observedGeneration:
              description: The generation observed by the deployment controller.
              format: int64
              type: integer
            readyReplicas:
              description: Total number of ready machines targeted by this deployment.
              format: int32
              type: integer
            replicas:
              description: Total number of non-terminated machines targeted by this
                deployment (their labels match the selector).
              format: int32
              type: integer
            unavailableReplicas:
              description: Total number of unavailable machines targeted by this deployment.
                This is the total number of machines that are still required for the
                deployment to have 100% available capacity. They may either be machines
                that are running but not yet available or machines that still have
                not been created.
              format: int32
              type: integer
            updatedReplicas:
              description: Total number of non-terminated machines targeted by this
                deployment that have the desired template spec.
              format: int32
              type: integer
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: machines.cluster.k8s.io
spec:
  group: cluster.k8s.io
  names:
    kind: Machine
    plural: machines
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            configSource:
              description: ConfigSource is used to populate in the associated Node
                for dynamic kubelet config. This field already exists in Node, so
                any updates to it in the Machine spec will be automatically copied
                to the linked NodeRef from the status. The rest of dynamic kubelet
                config support should then work as-is.
              type: object
            metadata:
              description: ObjectMeta will autopopulate the Node created. Use this
                to indicate what labels, annotations, name prefix, etc., should be
                used when creating the Node.
              type: object
            providerSpec:
              description: ProviderSpec details Provider-specific configuration to
                use during node creation.
              properties:
                value:
                  description: Value is an inlined, serialized representation of the
                    resource configuration. It is recommended that providers maintain
                    their own versioned API types that should be serialized/deserialized
                    from this field, akin to component config.
                  type: object
                valueFrom:
                  description: Source for the provider configuration. Cannot be used
                    if value is not empty.
                  properties:
                    machineClass:
                      description: The machine class from which the provider config
                        should be sourced.
                      properties:
                        provider:
                          description: Provider is the name of the cloud-provider
                            which MachineClass is intended for.
                          type: string
                      type: object
                  type: object
              type: object
            taints:
              description: Taints is the full, authoritative list of taints to apply
                to the corresponding Node. This list will overwrite any modifications
                made to the Node on an ongoing basis.
              items:
                type: object
              type: array
            versions:
              description: Versions of key software to use. This field is optional
                at cluster creation time, and omitting the field indicates that the
                cluster installation tool should select defaults for the user. These
                defaults may differ based on the cluster installer, but the tool should
                populate the values it uses when persisting Machine objects. A Machine
                spec missing this field at runtime is invalid.
              properties:
                controlPlane:
                  description: ControlPlane is the semantic version of the Kubernetes
                    control plane to run. This should only be populated when the machine
                    is a control plane.
                  type: string
                kubelet:
                  description: Kubelet is the semantic version of kubelet to run
                  type: string
              required:
              - kubelet
              type: object
          required:
          - providerSpec
          type: object
        status:
          properties:
            addresses:
              description: Addresses is a list of addresses assigned to the machine.
                Queried from cloud provider, if available.
              items:
                type: object
              type: array
            conditions:
              description: 'Conditions lists the conditions synced from the node conditions
                of the corresponding node-object. Machine-controller is responsible
                for keeping conditions up-to-date. MachineSet controller will be taking
                these conditions as a signal to decide if machine is healthy or needs
                to be replaced. Refer: https://kubernetes.io/docs/concepts/architecture/nodes/#condition'
              items:
                type: object
              type: array
            errorMessage:
              description: ErrorMessage will be set in the event that there is a terminal
                problem reconciling the Machine and will contain a more verbose string
                suitable for logging and human consumption.  This field should not
                be set for transitive errors that a controller faces that are expected
                to be fixed automatically over time (like service outages), but instead
                indicate that something is fundamentally wrong with the Machine's
                spec or the configuration of the controller, and that manual intervention
                is required. Examples of terminal errors would be invalid combinations
                of settings in the spec, values that are unsupported by the controller,
                or the responsible controller itself being critically misconfigured.  Any
                transient errors that occur during the reconciliation of Machines
                can be added as events to the Machine object and/or logged in the
                controller's output.
              type: string
            errorReason:
              description: ErrorReason will be set in the event that there is a terminal
                problem reconciling the Machine and will contain a succinct value
                suitable for machine interpretation.  This field should not be set
                for transitive errors that a controller faces that are expected to
                be fixed automatically over time (like service outages), but instead
                indicate that something is fundamentally wrong with the Machine's
                spec or the configuration of the controller, and that manual intervention
                is required. Examples of terminal errors would be invalid combinations
                of settings in the spec, values that are unsupported by the controller,
                or the responsible controller itself being critically misconfigured.  Any
                transient errors that occur during the reconciliation of Machines
                can be added as events to the Machine object and/or logged in the
                controller's output.
              type: string
            lastOperation:
              description: LastOperation describes the last-operation performed by
                the machine-controller. This API should be useful as a history in
                terms of the latest operation performed on the specific machine. It
                should also convey the state of the latest-operation for example if
                it is still on-going, failed or completed successfully.
              properties:
                description:
                  description: Description is the human-readable description of the
                    last operation.
                  type: string
                lastUpdated:
                  description: LastUpdated is the timestamp at which LastOperation
                    API was last-updated.
                  format: date-time
                  type: string
                state:
                  description: State is the current status of the last performed operation.
                    E.g. Processing, Failed, Successful etc
                  type: string
                type:
                  description: Type is the type of operation which was last performed.
                    E.g. Create, Delete, Update etc
                  type: string
              type: object
            lastUpdated:
              description: LastUpdated identifies when this status was last observed.
              format: date-time
              type: string
            nodeRef:
              description: NodeRef will point to the corresponding Node if it exists.
              type: object
            phase:
              description: Phase represents the current phase of machine actuation.
                E.g. Pending, Running, Terminating, Failed etc.
              type: string
            providerStatus:
              description: ProviderStatus details a Provider-specific status. It is
                recommended that providers maintain their own versioned API types
                that should be serialized/deserialized from this field.
              type: object
            versions:
              description: 'Versions specifies the current versions of software on
                the corresponding Node (if it exists). This is provided for a few
                reasons:  1) It is more convenient than checking the NodeRef, traversing
                it to    the Node, and finding the appropriate field in Node.Status.NodeInfo    (which
                uses different field names and formatting). 2) It removes some of
                the dependency on the structure of the Node,    so that if the structure
                of Node.Status.NodeInfo changes, only    machine controllers need
                to be updated, rather than every client    of the Machines API. 3)
                There is no other simple way to check the control plane    version.
                A client would have to connect directly to the apiserver    running
                on the target node in order to find out its version.'
              properties:
                controlPlane:
                  description: ControlPlane is the semantic version of the Kubernetes
                    control plane to run. This should only be populated when the machine
                    is a control plane.
                  type: string
                kubelet:
                  description: Kubelet is the semantic version of kubelet to run
                  type: string
              required:
              - kubelet
              type: object
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: machinesets.cluster.k8s.io
spec:
  group: cluster.k8s.io
  names:
    kind: MachineSet
    plural: machinesets
  scope: Namespaced
  subresources:
    scale:
      labelSelectorPath: .status.labelSelector
      specReplicasPath: .spec.replicas
      statusReplicasPath: .status.replicas
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            minReadySeconds:
              description: MinReadySeconds is the minimum number of seconds for which
                a newly created machine should be ready. Defaults to 0 (machine will
                be considered available as soon as it is ready)
              format: int32
              type: integer
            replicas:
              description: Replicas is the number of desired replicas. This is a pointer
                to distinguish between explicit zero and unspecified. Defaults to
                1.
              format: int32
              type: integer
            selector:
              description: 'Selector is a label query over machines that should match
                the replica count. Label keys and values that must match in order
                to be controlled by this MachineSet. It must match the machine template''s
                labels. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors'
              type: object
            template:
              description: Template is the object that describes the machine that
                will be created if insufficient replicas are detected.
              properties:
                metadata:
                  description: 'Standard object''s metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata'
                  type: object
                spec:
                  description: 'Specification of the desired behavior of the machine.
                    More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status'
                  properties:
                    configSource:
                      description: ConfigSource is used to populate in the associated
                        Node for dynamic kubelet config. This field already exists
                        in Node, so any updates to it in the Machine spec will be
                        automatically copied to the linked NodeRef from the status.
                        The rest of dynamic kubelet config support should then work
                        as-is.
                      type: object
                    metadata:
                      description: ObjectMeta will autopopulate the Node created.
                        Use this to indicate what labels, annotations, name prefix,
                        etc., should be used when creating the Node.
                      type: object
                    providerSpec:
                      description: ProviderSpec details Provider-specific configuration
                        to use during node creation.
                      properties:
                        value:
                          description: Value is an inlined, serialized representation
                            of the resource configuration. It is recommended that
                            providers maintain their own versioned API types that
                            should be serialized/deserialized from this field, akin
                            to component config.
                          type: object
                        valueFrom:
                          description: Source for the provider configuration. Cannot
                            be used if value is not empty.
                          properties:
                            machineClass:
                              description: The machine class from which the provider
                                config should be sourced.
                              properties:
                                provider:
                                  description: Provider is the name of the cloud-provider
                                    which MachineClass is intended for.
                                  type: string
                              type: object
                          type: object
                      type: object
                    taints:
                      description: Taints is the full, authoritative list of taints
                        to apply to the corresponding Node. This list will overwrite
                        any modifications made to the Node on an ongoing basis.
                      items:
                        type: object
                      type: array
                    versions:
                      description: Versions of key software to use. This field is
                        optional at cluster creation time, and omitting the field
                        indicates that the cluster installation tool should select
                        defaults for the user. These defaults may differ based on
                        the cluster installer, but the tool should populate the values
                        it uses when persisting Machine objects. A Machine spec missing
                        this field at runtime is invalid.
                      properties:
                        controlPlane:
                          description: ControlPlane is the semantic version of the
                            Kubernetes control plane to run. This should only be populated
                            when the machine is a control plane.
                          type: string
                        kubelet:
                          description: Kubelet is the semantic version of kubelet
                            to run
                          type: string
                      required:
                      - kubelet
                      type: object
                  required:
                  - providerSpec
                  type: object
              type: object
          required:
          - selector
          type: object
        status:
          properties:
            availableReplicas:
              description: The number of available replicas (ready for at least minReadySeconds)
                for this MachineSet.
              format: int32
              type: integer
            errorMessage:
              type: string
            errorReason:
              description: In the event that there is a terminal problem reconciling
                the replicas, both ErrorReason and ErrorMessage will be set. ErrorReason
                will be populated with a succinct value suitable for machine interpretation,
                while ErrorMessage will contain a more verbose string suitable for
                logging and human consumption.  These fields should not be set for
                transitive errors that a controller faces that are expected to be
                fixed automatically over time (like service outages), but instead
                indicate that something is fundamentally wrong with the MachineTemplate's
                spec or the configuration of the machine controller, and that manual
                intervention is required. Examples of terminal errors would be invalid
                combinations of settings in the spec, values that are unsupported
                by the machine controller, or the responsible machine controller itself
                being critically misconfigured.  Any transient errors that occur during
                the reconciliation of Machines can be added as events to the MachineSet
                object and/or logged in the controller's output.
              type: string
            fullyLabeledReplicas:
              description: The number of replicas that have labels matching the labels
                of the machine template of the MachineSet.
              format: int32
              type: integer
            observedGeneration:
              description: ObservedGeneration reflects the generation of the most
                recently observed MachineSet.
              format: int64
              type: integer
            readyReplicas:
              description: The number of ready replicas for this MachineSet. A machine
                is considered ready when the node has been created and is "Ready".
              format: int32
              type: integer
            replicas:
              description: Replicas is the most recently observed number of replicas.
              format: int32
              type: integer
          required:
          - replicas
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: cluster-api-manager-role
rules:
- apiGroups:
  - cluster.k8s.io
  resources:
  - clusters
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - machines
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - machinedeployments
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - machinesets
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - machines
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - cluster.k8s.io
  resources:
  - machines
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: null
  name: cluster-api-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-api-manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: cluster-api-system
---
apiVersion: v1
kind: Service
metadata:
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: cluster-api-controller-manager-service
  namespace: cluster-api-system
spec:
  ports:
  - port: 443
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: cluster-api-controller-manager
  namespace: cluster-api-system
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
      controller-tools.k8s.io: "1.0"
  serviceName: cluster-api-controller-manager-service
  template:
    metadata:
      labels:
        control-plane: controller-manager
        controller-tools.k8s.io: "1.0"
    spec:
      containers:
      - command:
        - /manager
        image: gcr.io/k8s-cluster-api/cluster-api-controller:latest
        name: manager
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
      terminationGracePeriodSeconds: 10
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoExecute
        key: node.alpha.kubernetes.io/notReady
        operator: Exists
      - effect: NoExecute
        key: node.alpha.kubernetes.io/unreachable
        operator: Exists
---
`
