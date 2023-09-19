import { FC } from 'react';

export interface ProfileCardProps {
  picture?: string;
  owner_name?: string;
  cardSize?: 'small' | 'medium' | 'large';
}

const ProfileCard: FC<ProfileCardProps> = ({ picture, owner_name, cardSize = 'small' }) => {
  return (
    <div className="flex gap-1.5">
      <img
        src={picture}
        alt=""
        className={
          cardSize === 'small'
            ? 'h-6 w-6 rounded-full'
            : cardSize === 'medium'
            ? 'h-8 w-8 rounded-full'
            : cardSize === 'large'
            ? 'h-10 w-10 rounded-full'
            : ''
        }
      />
      <div className="mt-0.5" data-cy="user-email">
        {owner_name}
      </div>
    </div>
  );
};

export { ProfileCard };
