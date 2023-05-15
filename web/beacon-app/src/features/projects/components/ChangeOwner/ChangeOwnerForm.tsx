import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikProvider } from 'formik';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import Select from '@/components/ui/Select';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import type { Member } from '@/features/teams/types/member';
import { MemberStatusEnum } from '@/features/teams/types/member';

import type { Project } from '../../types/Project';
import {
  UpdateProjectOwnerFormDTO,
  useUpdateProjectOwnerForm,
} from '../../types/updateProjectOwnerFormService';
type ChangeOwnerFormProps = {
  handleSubmit: (values: UpdateProjectOwnerFormDTO) => void;
  initialValues: Project;
};

const ChangeOwnerForm = ({ handleSubmit, initialValues }: ChangeOwnerFormProps) => {
  const formik = useUpdateProjectOwnerForm(handleSubmit, initialValues);
  const { getFieldProps, isSubmitting, values, setFieldValue } = formik;
  const { members } = useFetchMembers();

  const formatMembers = () => {
    return members?.members?.map((member: Member) => ({
      label: member?.name,
      value: member?.id,
      status: member?.status,
    }));
  };

  const optionsAvailable = () =>
    formatMembers()?.filter(
      (opt: any) =>
        opt?.value !== values.current_owner.value && opt?.status !== MemberStatusEnum.PENDING
    );

  useEffect(() => {
    if (optionsAvailable()?.length === 0) {
      toast.error(t`There are no other members to select as the new owner.`);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [values.current_owner.value, optionsAvailable()]);

  const getDefaultOption = () => {
    if (optionsAvailable?.length === 0) {
      return null;
    }
    // select the first option
    return optionsAvailable()[0];
  };

  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3" data-testid="update-owner-form">
        <TextField
          label={t`Current Owner`}
          {...getFieldProps('current_owner.label')}
          isDisabled
          data-testid="current-owner"
          data-cy="prj-current-owner"
        />
        <fieldset>
          <label htmlFor="new_owner" className="text-sm">
            <Trans>Select New Owner</Trans>
          </label>
          <Select
            id="new_owner"
            inputId="new_owner"
            isDisabled={isSubmitting}
            defaultValue={getDefaultOption}
            options={optionsAvailable()}
            name="new_owner"
            onChange={(value: any) => setFieldValue('new_owner', value.value)}
          />
        </fieldset>
        <ErrorMessage name="new_owner" component="small" className="text-xs text-danger-500" />
        <div className="pt-3 text-center">
          <Button
            type="submit"
            isLoading={isSubmitting}
            disabled={isSubmitting || !values?.new_owner}
            data-cy="update-owner"
            data-testid="update-owner"
          >
            <Trans>Save</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default ChangeOwnerForm;
