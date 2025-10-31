import React from 'react';
import { Role } from '../types';
import { cn } from '../utils/cn';

interface RoleSelectorProps {
  roles: Role[];
  selectedRole: Role | null;
  onSelectRole: (role: Role) => void;
  isLoading?: boolean;
}

export const RoleSelector: React.FC<RoleSelectorProps> = ({
  roles,
  selectedRole,
  onSelectRole,
  isLoading = false,
}) => {
  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold text-gray-800 text-center">
        选择你的AI角色
      </h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {roles.map((role) => (
          <div
            key={role.id}
            className={cn(
              "card p-6 cursor-pointer transition-all duration-200 hover:shadow-xl hover:scale-105",
              selectedRole?.id === role.id
                ? "ring-2 ring-primary-500 bg-primary-50"
                : "hover:bg-gray-50"
            )}
            onClick={() => onSelectRole(role)}
          >
            <div className="text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-primary-400 to-primary-600 rounded-full flex items-center justify-center text-white text-2xl font-bold">
                {role.name.charAt(0)}
              </div>
              <h3 className="text-lg font-semibold text-gray-800 mb-2">
                {role.name}
              </h3>
              <div className="flex flex-wrap gap-1 justify-center">
                {role.description}
              </div>
            </div>
          </div>
        ))}
      </div>
      {isLoading && (
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <p className="mt-2 text-gray-600">加载角色中...</p>
        </div>
      )}
    </div>
  );
};
