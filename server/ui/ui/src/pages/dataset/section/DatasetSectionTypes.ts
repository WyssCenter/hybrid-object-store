

export interface Group {
  group_name: string;
}

export interface User {
  user_name: string;
}


export interface FetchPermissionItem {
  group?: Group;
  user?:User;
  permission: string;
}



export interface PermissionItem {
  name: string;
}

export interface Permission {
  group?: Group;
  user?: User;
  permission: PermissionItem;
}

export interface Dataset {
  description: string;
  permissions: Array<FetchPermissionItem>;
  sync_type: string;
  sync_policy: string;
  sync_enabled: boolean;
}


export interface GroupItem {
  group_name: string;
}
