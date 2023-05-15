import { FC } from 'react';

export interface ProfileCardProps {
  picture?: string;
  owner_name?: string;
}

const ProfileCard: FC<ProfileCardProps> = ({ picture, owner_name }) => {
  return (
    <div className="flex gap-1.5">
      <img src={picture} alt="" className="h-6 w-6 rounded-2xl" />
      <div className="mt-0.5">{owner_name}</div>
    </div>
  );
};

export { ProfileCard };
