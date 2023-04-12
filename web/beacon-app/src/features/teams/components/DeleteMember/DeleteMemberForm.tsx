import { Checkbox } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikProvider } from 'formik';
import styled from 'styled-components';

import Button from '@/components/ui/Button';
import TextField from '@/components/ui/TextField';
import {
  DeleteMemberFormValue,
  useDeleteMemberForm,
} from '@/features/members/types/deleteMemberFormService';

type NewMemberFormProps = {
  onSubmit: (values: DeleteMemberFormValue) => void;
  isSubmitting?: boolean;
  initialValues: DeleteMemberFormValue;
};

function DeleteMemberForm({ onSubmit, isSubmitting, initialValues }: NewMemberFormProps) {
  const formik = useDeleteMemberForm(onSubmit, initialValues);

  const { getFieldProps, setFieldValue, values } = formik;

  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField type="hidden" {...getFieldProps('id')} />
        <StyledTextField
          label="Team Member"
          {...getFieldProps('name')}
          isDisabled
          data-testid="name"
        />
        <CheckboxFieldset>
          <Checkbox
            name="delete_agreement"
            onChange={(isSelected) => {
              setFieldValue('delete_agreement', isSelected);
            }}
            data-testid="delete_agreement"
          >
            Check to confirm removal. The team member will no longer have access to the
            organization.
          </Checkbox>
          <ErrorMessage name="delete_agreement" component="p" />
        </CheckboxFieldset>
        <div className="pt-3 text-center">
          <Button
            type="submit"
            isLoading={isSubmitting}
            isDisabled={!values.delete_agreement}
            data-testid="remove-btn"
          >
            Remove
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default DeleteMemberForm;

const CheckboxFieldset = styled.fieldset`
  margin-top: 1rem;
  font-weight: 500;
  font-size: 0.85rem;
  width: 40vh;
  label svg {
    min-width: 30px;
  }
`;

const StyledTextField = styled(TextField)`
  width: 40vh;
`;
