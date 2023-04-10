import { Form, FormikHelpers, FormikProvider } from 'formik';

import Button from '@/components/ui/Button';
import Select from '@/components/ui/Select';
import TextField from '@/components/ui/TextField';
import { ROLE_OPTIONS, useNewMemberForm } from '@/features/members/types/memberFormService';
import type { NewMemberDTO } from '@/features/members/types/memberServices';

type NewMemberFormProps = {
  onSubmit: (values: NewMemberDTO, helpers: FormikHelpers<NewMemberDTO>) => void;
  isDisabled?: boolean;
  isLoading?: boolean;
};

function NewMemberForm({ onSubmit, isDisabled, isLoading }: NewMemberFormProps) {
  const formik = useNewMemberForm(onSubmit);

  const { touched, errors, getFieldProps, setFieldValue, values } = formik;
  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField
          label="Team Member"
          placeholder="Natali Craig"
          errorMessage={touched.email && errors.email}
          {...getFieldProps('name')}
          isDisabled
        />
        <TextField
          label="Current role"
          placeholder="Member"
          {...getFieldProps('current_role')}
          isDisabled
        />
        <fieldset>
          <label htmlFor="role" className="text-sm">
            Select new role
          </label>
          <Select
            id="role"
            isDisabled={isLoading}
            defaultValue={ROLE_OPTIONS.filter((opt) => opt.value === values.current_role)}
            options={ROLE_OPTIONS.filter((opt) => opt.value !== values.current_role)}
            name="role"
            value={ROLE_OPTIONS.filter((opt) => opt.value === values.role)}
            onChange={(value: any) => setFieldValue('role', value.value)}
          />
        </fieldset>
        <div className="pt-3 text-center">
          <Button type="submit" isLoading={isLoading} isDisabled={isDisabled}>
            Save
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default NewMemberForm;
