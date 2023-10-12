import { getInitials } from '@/utils/strings';

type ProfileAvatar = {
  name: string;
};

const ProfileAvatar = ({ name }: ProfileAvatar) => {
  return (
    <div className="flex h-[28px] w-9 items-center justify-center rounded-full bg-[#F7F9FB] md:bg-primary">
      <span className="text-sm font-semibold md:text-white">{getInitials(name)}</span>
    </div>
  );
};

export default ProfileAvatar;
