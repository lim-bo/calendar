#include "event_entry.h"
#include "ui_event_entry.h"

event_entry::event_entry(EventData data, QString viewerUID,  QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::event_entry), data(data), viewerUID(viewerUID)
{
    ui->setupUi(this);
    ui->name->setText(data.name);
    ui->type->setText(data.type);
    QString priorityText;
    switch(data.prior) {
    case 1: priorityText = "Менее важно"; break;
    case 2: priorityText = "Средне важно"; break;
    case 3: priorityText = "Важно"; break;
    }
    ui->prior->setText(priorityText);
    QString startStr = data.start.toString("dd.MM.yyyy hh:mm");
    QString endStr = data.end.toString("dd.MM.yyyy hh:mm");
    ui->vremyas->setText(startStr + " — " + endStr);
    ui->desc->setText(data.desc);

    ui->checkBox->setMinimumSize(100, 20);

    bool isParticipant = false;
    for (const Participant& participant : data.parts) {
        if (participant.uid == viewerUID) {
            isParticipant = participant.accepted;
            break;
        }
    }

    if (data.master == viewerUID) {
        ui->status->setText("Вы создатель");
    } else {
        ui->status->setText(isParticipant ? "Вас добавили участником" : "Вы не участник");
    }

    isAccepted = false;

    ui->checkBox->setChecked(isAccepted);
    ui->checkBox->setText(isAccepted ? "Участие подтверждено" : "Подтвердждение участия");


    connect(ui->checkBox, &QCheckBox::stateChanged, this, &event_entry::on_checkBox_stateChanged);
}

const EventData event_entry::getData() const {
    return this->data;
}

event_entry::~event_entry()
{
    delete ui;
}

void event_entry::on_pushButton_clicked()
{
    emit deleted(this);
}


void event_entry::on_pushButton_3_clicked()
{
    ChatWindow chat(data.ID, viewerUID);
    chat.exec();
}

void event_entry::on_editButton_clicked()
{
    Event* editWindow = new Event(viewerUID);
    editWindow->setAttribute(Qt::WA_DeleteOnClose);
    editWindow->fillFormWithEventData(this->data);
    connect(editWindow, &Event::eventUpdated, this, [this]() {
        emit edited(this);
    });
    editWindow->show();


}

void event_entry::on_pushButton_4_clicked()
{
    Client cli(settings.value("host").toString(), settings.value("port").toString());
    attachments window(this->data.ID, viewerUID, &cli, this);
    window.exec();
}


void event_entry::on_checkBox_stateChanged(int arg1)
{
    bool newState = (arg1 == Qt::Checked);
    if (newState != isAccepted) {
        Client cli(settings.value("host").toString(), settings.value("port").toString());
        bool success = cli.updateParticipation(data.ID, viewerUID, newState);

        if (success) {
            isAccepted = newState;
            ui->checkBox->setText(isAccepted ? "Участие подтверждено" : "Подтверждение присутствия");

            for (Participant& participant : data.parts) {
                if (participant.uid == viewerUID) {
                    participant.accepted = newState;
                    break;
                }
            }
        } else {

            ui->checkBox->blockSignals(true);
            ui->checkBox->setChecked(isAccepted);
            ui->checkBox->blockSignals(false);

        }
    }
}


void event_entry::on_part_clicked()
{
    allparticipants dialog(data, this);
    dialog.exec();
}


void event_entry::on_pushButton_5_clicked()
{
    replay_event dialog(this->data);

    if (dialog.exec() == QDialog::Accepted) {
        EventData repeatedEvent = dialog.getRepeatedEvent();

        Client cli(settings.value("host").toString(), settings.value("port").toString());
        QString newEventId;
        bool success = cli.addEvent(repeatedEvent, viewerUID, &newEventId);

        if (success) {
            QMessageBox::information(this, "Успех", "Событие успешно повторено");
            emit edited(this);
        } else {
            QMessageBox::warning(this, "Ошибка", "Не удалось повторить событие");
        }
    }
}

