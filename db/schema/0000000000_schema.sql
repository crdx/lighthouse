create table sessions (
    k varchar(64) not null default '',
    v blob not null,
    e bigint(20) not null default 0,
    primary key (k)
);

create table devices (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    name varchar(100) not null,
    hostname varchar(100) not null,
    state varchar(20) not null,
    last_seen_at datetime not null,
    grace_period text not null,
    icon varchar(100) not null,
    notes text not null,
    watch tinyint(1) not null,
    hostname_announced_at datetime default null,
    origin tinyint(1) not null,
    state_updated_at datetime not null,
    `limit` text not null,
    ping tinyint(1) not null,
    primary key (id),
    key idx_devices_deleted_at (deleted_at),
    constraint chk_devices__icon_valid check (icon regexp '^.+:.+$' or icon = ''),
    constraint chk_devices__watch_valid check (watch in (0,1)),
    constraint chk_devices__origin_valid check (origin in (0,1)),
    constraint chk_devices__grace_period_valid check (grace_period <> ''),
    constraint chk_devices__state_valid check (state in ('online','offline'))
);

create table users (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    username varchar(20) not null,
    password_hash text not null,
    last_login datetime default null,
    role bigint not null,
    last_login_at datetime default null,
    last_visit_at datetime default null,
    primary key (id),
    key idx_users_deleted_at (deleted_at),
    constraint chk_users__role_valid check (role in (1,2,3)),
    constraint chk_users__last_login_at_valid check (last_login_at >= created_at)
);

create table adapters (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    name varchar(100) not null,
    mac_address varchar(17) not null,
    vendor varchar(100) not null,
    ip_address varchar(15) not null,
    last_seen_at datetime not null,
    primary key (id),
    key idx_adapters_deleted_at (deleted_at),
    key idx_adapters_mac_address (mac_address),
    key fk_adapters__device (device_id),
    constraint fk_adapters__device foreign key (device_id) references devices (id),
    constraint chk_adapters__mac_address_valid check (mac_address regexp '^[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$'),
    constraint chk_adapters__ip_address_valid check (ip_address regexp '^[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}$')
);

create table audit_logs (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    user_id bigint default null,
    ip_address varchar(15) not null,
    message text not null,
    device_id bigint default null,
    primary key (id),
    key idx_audit_logs_deleted_at (deleted_at),
    key fk_audit_logs__user (user_id),
    key fk_audit_logs__device (device_id),
    constraint fk_audit_logs__device foreign key (device_id) references devices (id),
    constraint fk_audit_logs__user foreign key (user_id) references users (id),
    constraint chk_audit_logs__ip_address_valid check (ip_address regexp '^[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}$'),
    constraint chk_audit_logs__message_valid check (message <> '')
);

create table device_discovery_notifications (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    processed tinyint(1) not null DEFAULT 0,
    primary key (id),
    key idx_device_discovery_notifications_deleted_at (deleted_at),
    key idx_device_discovery_notifications_device_id (device_id),
    constraint fk_device_discovery_notifications__device foreign key (device_id) references devices (id),
    constraint chk_device_discovery_notifications__processed_valid check (processed in (0,1))
);

create table device_ip_address_logs (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    ip_address text not null,
    primary key (id),
    key idx_device_ip_address_logs_device_id (device_id),
    key idx_device_ip_address_logs_deleted_at (deleted_at),
    constraint fk_device_ip_address_logs__device foreign key (device_id) references devices (id),
    constraint chk_device_ip_address_logs__ip_address_valid check (ip_address regexp '^[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}$')
);

create table device_limit_notifications (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    state_updated_at datetime not null,
    `limit` text not null,
    processed tinyint(1) not null,
    primary key (id),
    key idx_device_limit_notifications_deleted_at (deleted_at),
    key idx_device_limit_notifications_device_id (device_id),
    constraint fk_device_limit_notifications__device foreign key (device_id) references devices (id),
    constraint chk_device_limit_notifications__processed_valid check (processed in (0,1))
);

create table device_state_logs (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    state varchar(20) not null,
    grace_period text not null,
    primary key (id),
    key idx_device_state_logs_deleted_at (deleted_at),
    key fk_device_state_logs__device (device_id),
    constraint fk_device_state_logs__device foreign key (device_id) references devices (id),
    constraint chk_device_state_logs__state_valid check (state in ('online','offline'))
);

create table device_state_notifications (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    state varchar(20) not null,
    processed tinyint(1) not null,
    grace_period text not null,
    primary key (id),
    key idx_device_state_notifications_device_id (device_id),
    key idx_device_state_notifications_deleted_at (deleted_at),
    constraint fk_device_state_notifications__device foreign key (device_id) references devices (id),
    constraint chk_device_state_notifications__state_valid check (state in ('online','offline')),
    constraint chk_device_state_notifications__processed_valid check (processed in (0,1))
);

create table mappings (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    mac_address varchar(17) not null,
    ip_address varchar(15) not null,
    label varchar(20) not null,
    primary key (id),
    key idx_mappings_mac_address (mac_address),
    key idx_mappings_deleted_at (deleted_at),
    constraint chk_mappings__mac_address_valid check (mac_address regexp '^[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$'),
    constraint chk_mappings__ip_address_valid check (ip_address regexp '^[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}.[0-9]{0,3}$')
);

create table notifications (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    subject varchar(200) not null,
    body text not null,
    primary key (id),
    key idx_notifications_deleted_at (deleted_at)
);

create table scan_results (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    scan_id bigint not null,
    device_id bigint not null,
    port bigint not null,
    primary key (id),
    key idx_scan_results_deleted_at (deleted_at)
);

create table scans (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    completed_at datetime default null,
    primary key (id),
    key idx_scans_deleted_at (deleted_at)
);

create table services (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    port bigint not null,
    name text not null,
    last_seen_at datetime not null,
    primary key (id),
    key idx_services_deleted_at (deleted_at),
    key fk_services__device (device_id),
    constraint fk_services__device foreign key (device_id) references devices (id)
);

create table device_service_notifications (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    device_id bigint not null,
    service_id bigint default null,
    processed tinyint(1) not null DEFAULT 0,
    primary key (id),
    key idx_device_service_notifications_device_id (device_id),
    key idx_device_service_notifications_deleted_at (deleted_at),
    key fk_device_service_notifications__service (service_id),
    constraint fk_device_service_notifications__service foreign key (service_id) references services (id)
);

create table settings (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    name varchar(100) not null,
    value text not null,
    primary key (id),
    key idx_settings_deleted_at (deleted_at)
);

create table vendor_lookups (
    id bigint not null auto_increment,
    created_at datetime not null,
    updated_at datetime default null,
    deleted_at datetime default null,
    adapter_id bigint not null,
    processed tinyint(1) not null,
    succeeded tinyint(1) not null,
    primary key (id),
    key idx_vendor_lookups_deleted_at (deleted_at),
    key idx_vendor_lookups_processed (processed),
    key fk_vendor_lookups__adapter (adapter_id),
    constraint fk_vendor_lookups__adapter foreign key (adapter_id) references adapters (id),
    constraint chk_vendor_lookups__processed_valid check (processed in (0,1)),
    constraint chk_vendor_lookups__succeeded_valid check (succeeded in (0,1))
);
