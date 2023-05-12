import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikProvider } from 'formik';
import { useCallback } from 'react';

import Select from '@/components/ui/Select';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import type { Member } from '@/features/teams/types/member';

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
  console.log('[] memebers data', members);
  const formatMembers = useCallback(() => {
    return members?.members?.map((member: Member) => ({
      label: member?.name,
      value: member?.id,
    }));
  }, [members]);

  console.log('[] formatMembers', formatMembers());
  console.log('[] values', values);
  console.log('[] initialValues', initialValues);

  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField
          label={t`Current owner`}
          {...getFieldProps('current_owner.name')}
          isDisabled
          data-cy="prj-current-owner"
        />
        <fieldset>
          <label htmlFor="role" className="text-sm">
            <Trans>Select New Owner</Trans>
          </label>
          <Select
            id="role"
            isDisabled={isSubmitting}
            defaultValue={formatMembers()?.filter(
              (opt: any) => opt?.value === values.current_owner.name
            )}
            options={formatMembers()?.filter((opt: any) => opt?.value !== values.current_owner?.id)}
            name="new_owner"
            value={formatMembers()?.filter((opt: any) => opt?.value === values.current_owner?.id)}
            onChange={(value: any) => setFieldValue('new_owner', value.value)}
          />
          <ErrorMessage name="new_owner" component="small" className="text-xs text-danger-500" />
        </fieldset>
        <div className="pt-3 text-center">
          <Button
            type="submit"
            isLoading={isSubmitting}
            disabled={isSubmitting}
            data-cy="saveNewRole"
          >
            <Trans>Save</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default ChangeOwnerForm;
