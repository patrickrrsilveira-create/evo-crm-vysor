# frozen_string_literal: true

# Backfills new roles management permissions
# (roles.create, roles.update, roles.delete, permissions.create, permissions.update, permissions.delete)
# to already-bootstrapped installations.
class AddRolesManagementToAccountOwner < ActiveRecord::Migration[7.1]
  PERMISSIONS = %w[
    roles.create
    roles.update
    roles.delete
    roles.bulk_assign
    roles.bulk_update_permissions
    permissions.create
    permissions.update
    permissions.delete
    permissions.assign
    permissions.bulk_operations
  ].freeze

  ROLE_KEYS = %w[account_owner super_admin].freeze

  def up
    return unless ActiveRecord::Base.connection.table_exists?(:roles)

    ROLE_KEYS.each do |role_key|
      role = Role.find_by(key: role_key)
      next unless role

      PERMISSIONS.each do |permission_key|
        next if role.role_permissions_actions.exists?(permission_key: permission_key)

        # Ignore if it is not valid in ResourceActionsConfig
        next unless ResourceActionsConfig.valid_permission?(permission_key)

        role.role_permissions_actions.create!(permission_key: permission_key)
      end
    end
  end

  def down
    return unless ActiveRecord::Base.connection.table_exists?(:roles)

    ROLE_KEYS.each do |role_key|
      role = Role.find_by(key: role_key)
      next unless role

      role.role_permissions_actions.where(permission_key: PERMISSIONS).destroy_all
    end
  end
end
